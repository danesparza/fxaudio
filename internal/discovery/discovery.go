// Package discovery advertises the FX HTTP API for local fxcontroller discovery.
package discovery

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/grandcat/zeroconf"
	"github.com/rs/zerolog/log"
)

const ServiceType = "_fx._tcp"
const Domain = "local."

// Config is shared by all FX services. ID should be unique per installation.
type Config struct {
	Enabled bool
	Name    string
	ID      string
	Service string
}

func (c Config) records(host string, port int) (string, []string, error) {
	name := c.Name
	if name == "" {
		name = fmt.Sprintf("%s-%s-%d", c.Service, host, port)
	}
	if !utf8.ValidString(name) || strings.ContainsAny(name, "\x00\r\n") || strings.TrimSpace(name) == "" || len(name) > 63 {
		return "", nil, fmt.Errorf("discovery name must contain 1–63 UTF-8 bytes")
	}
	id := c.ID
	if id == "" {
		id = fmt.Sprintf("%s:%s:%d", c.Service, host, port)
	}
	records := []string{"txtvers=1", "service=" + c.Service, "id=" + id, "api=v1", "scheme=http", "path=/v1"}
	for _, record := range records {
		if len(record) > 255 || strings.ContainsAny(record, "\x00\r\n") || !utf8.ValidString(record) {
			return "", nil, fmt.Errorf("invalid discovery TXT record")
		}
	}
	return name, records, nil
}

func advertise(c Config, addr *net.TCPAddr) (*zeroconf.Server, error) {
	host, err := os.Hostname()
	if err != nil {
		return nil, err
	}
	name, records, err := c.records(host, addr.Port)
	if err != nil {
		return nil, err
	}
	if addr.IP.IsUnspecified() && addr.IP.To4() == nil {
		return zeroconf.Register(name, ServiceType, Domain, addr.Port, records, nil)
	}
	// A listener bound to one address must not advertise other host addresses.
	if addr.IP.IsLoopback() {
		return nil, fmt.Errorf("loopback listener is not discoverable on the LAN")
	}
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	// An IPv4-only wildcard listener must not publish unreachable IPv6 addresses.
	if addr.IP.IsUnspecified() {
		var ips []string
		var selected []net.Interface
		for _, iface := range interfaces {
			if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagMulticast == 0 || iface.Flags&net.FlagLoopback != 0 {
				continue
			}
			addresses, err := iface.Addrs()
			if err != nil {
				continue
			}
			count := len(ips)
			for _, address := range addresses {
				ip, _, err := net.ParseCIDR(address.String())
				if err == nil && ip.To4() != nil {
					ips = append(ips, ip.String())
				}
			}
			if len(ips) > count {
				selected = append(selected, iface)
			}
		}
		if len(ips) == 0 {
			return nil, fmt.Errorf("no multicast IPv4 interface available")
		}
		return zeroconf.RegisterProxy(name, ServiceType, Domain, addr.Port, host, ips, records, selected)
	}
	for _, iface := range interfaces {
		addresses, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, address := range addresses {
			ip, _, err := net.ParseCIDR(address.String())
			if err == nil && ip.Equal(addr.IP) {
				return zeroconf.RegisterProxy(name, ServiceType, Domain, addr.Port, host, []string{addr.IP.String()}, records, []net.Interface{iface})
			}
		}
	}
	return nil, fmt.Errorf("no network interface matches listener %s", addr)
}

// ListenAndServe binds HTTP before advertising it. Discovery failures are logged
// but do not prevent direct API access. Cancellation withdraws the advertisement
// and drains HTTP requests for up to five seconds.
func ListenAndServe(ctx context.Context, address string, handler http.Handler, config Config) error {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	return serve(ctx, listener, handler, config)
}

func serve(ctx context.Context, listener net.Listener, handler http.Handler, config Config) error {
	defer listener.Close()
	var announcement *zeroconf.Server
	if config.Enabled {
		var err error
		announcement, err = advertise(config, listener.Addr().(*net.TCPAddr))
		if err != nil {
			log.Warn().Err(err).Str("service", config.Service).Msg("Zeroconf advertisement unavailable; HTTP API remains available")
		} else {
			log.Info().Str("type", ServiceType).Str("service", config.Service).Str("address", listener.Addr().String()).Msg("Zeroconf advertisement started")
		}
	}
	server := &http.Server{Handler: handler, ReadHeaderTimeout: 10 * time.Second}
	done := make(chan struct{})
	stopped := make(chan struct{})
	go func() {
		defer close(stopped)
		select {
		case <-ctx.Done():
			if announcement != nil {
				announcement.Shutdown()
			}
			shutdownCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 5*time.Second)
			defer cancel()
			if err := server.Shutdown(shutdownCtx); err != nil {
				log.Warn().Err(err).Msg("HTTP shutdown timed out; closing active connections")
				_ = server.Close()
			}
		case <-done:
			if announcement != nil {
				announcement.Shutdown()
			}
		}
	}()
	log.Info().Str("address", listener.Addr().String()).Msg("HTTP API listener started")
	err := server.Serve(listener)
	close(done)
	<-stopped
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	return err
}
