package discovery

import (
	"context"
	"io"
	"net"
	"net/http"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestRecords(t *testing.T) {
	for _, service := range []string{"fxaudio", "fxpixel", "fxdmx", "fxtrigger"} {
		t.Run(service, func(t *testing.T) {
			for _, tc := range []struct {
				name             string
				config           Config
				wantName, wantID string
				invalid          bool
			}{
				{name: "defaults", config: Config{Service: service}, wantName: service + "-stage-3030", wantID: service + ":stage:3030"},
				{name: "custom", config: Config{Service: service, Name: "Front porch", ID: "installation-1"}, wantName: "Front porch", wantID: "installation-1"},
				{name: "long name", config: Config{Service: service, Name: strings.Repeat("a", 64)}, invalid: true},
				{name: "blank name", config: Config{Service: service, Name: " "}, invalid: true},
				{name: "long ID", config: Config{Service: service, ID: strings.Repeat("a", 253)}, invalid: true},
				{name: "invalid ID", config: Config{Service: service, ID: "a\x00b"}, invalid: true},
				{name: "UTF8 byte limit", config: Config{Service: service, Name: strings.Repeat("é", 32)}, invalid: true},
			} {
				t.Run(tc.name, func(t *testing.T) {
					name, records, err := tc.config.records("stage", 3030)
					if tc.invalid {
						if err == nil {
							t.Fatal("expected invalid metadata error")
						}
						return
					}
					if err != nil {
						t.Fatal(err)
					}
					want := []string{"txtvers=1", "service=" + service, "id=" + tc.wantID, "api=v1", "scheme=http", "path=/v1"}
					if name != tc.wantName || !reflect.DeepEqual(records, want) {
						t.Fatalf("got %q %v; want %q %v", name, records, tc.wantName, want)
					}
				})
			}
		})
	}
}

func TestServeLifecycle(t *testing.T) {
	for _, enabled := range []bool{false, true} {
		name := "disabled"
		if enabled {
			name = "unavailable on loopback"
		}
		t.Run(name, func(t *testing.T) {
			listener, err := net.Listen("tcp", "127.0.0.1:0")
			if err != nil {
				t.Fatal(err)
			}
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()
			done := make(chan error, 1)
			go func() {
				done <- serve(ctx, listener, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { _, _ = io.WriteString(w, "ready") }), Config{Enabled: enabled, Service: "fxaudio"})
			}()
			client := &http.Client{Timeout: time.Second, Transport: &http.Transport{Proxy: nil}}
			defer client.CloseIdleConnections()
			req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, "http://"+listener.Addr().String(), nil)
			if err != nil {
				t.Fatal(err)
			}
			resp, err := client.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			body, err := io.ReadAll(resp.Body)
			resp.Body.Close()
			if err != nil || string(body) != "ready" {
				t.Fatalf("response %q, error %v", body, err)
			}
			cancel()
			select {
			case err := <-done:
				if err != nil {
					t.Fatal(err)
				}
			case <-time.After(2 * time.Second):
				t.Fatal("server did not stop")
			}
			conn, err := net.DialTimeout("tcp", listener.Addr().String(), time.Second)
			if err == nil {
				conn.Close()
				t.Fatal("listener still accepts connections")
			}
		})
	}
}

func TestListenFailure(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()
	if err := ListenAndServe(t.Context(), listener.Addr().String(), http.NotFoundHandler(), Config{Enabled: true}); err == nil {
		t.Fatal("expected occupied port error")
	}
}
