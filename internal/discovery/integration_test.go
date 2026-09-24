package discovery

import (
	"context"
	"fmt"
	"net"
	"os"
	"testing"
	"time"

	"github.com/grandcat/zeroconf"
)

// Run explicitly on a multicast-capable LAN; ordinary unit tests need no mDNS.
func TestZeroconfIntegration(t *testing.T) {
	if os.Getenv("FX_DISCOVERY_INTEGRATION") != "1" {
		t.Skip("set FX_DISCOVERY_INTEGRATION=1 for live multicast verification")
	}
	for _, network := range []string{"tcp", "tcp4"} {
		t.Run(network, func(t *testing.T) { checkAdvertisement(t, network) })
	}
}

func checkAdvertisement(t *testing.T, network string) {
	listener, err := net.Listen(network, ":0")
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()
	name := fmt.Sprintf("fx-test-%d", time.Now().UnixNano())
	server, err := advertise(Config{Service: "fxaudio", Name: name, ID: name}, listener.Addr().(*net.TCPAddr))
	if err != nil {
		t.Fatal(err)
	}
	defer server.Shutdown()
	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	defer cancel()
	entries := make(chan *zeroconf.ServiceEntry)
	if err := resolver.Browse(ctx, ServiceType, Domain, entries); err != nil {
		t.Fatal(err)
	}
	for {
		select {
		case entry := <-entries:
			if entry == nil {
				t.Fatal("discovery channel closed")
			}
			if entry.Instance != name {
				continue
			}
			if entry.Port != listener.Addr().(*net.TCPAddr).Port {
				t.Fatalf("wrong port: %d", entry.Port)
			}
			if len(entry.AddrIPv4)+len(entry.AddrIPv6) == 0 {
				t.Fatal("no advertised addresses")
			}
			if network == "tcp4" && len(entry.AddrIPv6) != 0 {
				t.Fatal("IPv4-only listener advertised IPv6 addresses")
			}
			found := false
			for _, record := range entry.Text {
				if record == "id="+name {
					found = true
				}
			}
			if !found {
				t.Fatalf("missing identity: %v", entry.Text)
			}
			return
		case <-ctx.Done():
			t.Fatal("advertisement was not discovered on the LAN")
		}
	}
}
