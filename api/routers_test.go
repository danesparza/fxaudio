package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
)

func TestRouterPreservesPeerAddress(t *testing.T) {
	for _, tc := range []struct{ name, header, value string }{
		{name: "no forwarded header"},
		{name: "true client", header: "True-Client-IP", value: "198.51.100.2"},
		{name: "real IP", header: "X-Real-IP", value: "198.51.100.2"},
		{name: "forwarded for", header: "X-Forwarded-For", value: "198.51.100.2"},
		{name: "forwarded chain", header: "X-Forwarded-For", value: "198.51.100.2, 192.0.2.5"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/peer-address", nil)
			req.RemoteAddr = "192.0.2.1:1234"
			if tc.header != "" {
				req.Header.Set(tc.header, tc.value)
			}
			router := NewRouter(Service{}).(*chi.Mux)
			var peer string
			router.Get("/peer-address", func(w http.ResponseWriter, r *http.Request) { peer = r.RemoteAddr; w.WriteHeader(http.StatusNoContent) })
			response := httptest.NewRecorder()
			router.ServeHTTP(response, req)
			if response.Code != http.StatusNoContent {
				t.Fatalf("unexpected status: %d", response.Code)
			}
			if peer != "192.0.2.1:1234" {
				t.Fatalf("untrusted %s changed peer address to %q", tc.header, peer)
			}
		})
	}
}
