package api

import (
	"encoding/json"
	"fmt"
	"github.com/danesparza/fxaudio/internal/data"
	"github.com/danesparza/fxaudio/internal/media"
	"github.com/danesparza/fxaudio/version"
	"net/http"
	"time"
)

// Service encapsulates API service operations
type Service struct {
	DB        data.AppDataService
	StartTime time.Time

	// PlayMedia signals a file should be played
	PlayMedia chan media.PlayAudioRequest

	// StopMedia signals a file should stop playing
	StopMedia chan string

	//	StopAllMedia signals all audio should stop playing
	StopAllMedia chan bool
}

// PlayAudioRequest represents a request to play an audio endpoint
type PlayAudioRequest struct {
	Endpoint string `json:"endpoint"`
}

// SystemResponse is a response for a system request
type SystemResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// ErrorResponse represents an API response
type ErrorResponse struct {
	Message string `json:"message"`
}

// Used to send back an error:
func sendErrorResponse(rw http.ResponseWriter, err error, code int) {
	//	Our return value
	response := ErrorResponse{
		Message: "Error: " + err.Error()}

	//	Serialize to JSON & return the response:
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	rw.WriteHeader(code)
	json.NewEncoder(rw).Encode(response)
}

// ShowUI redirects to the /ui/ url path
func ShowUI(rw http.ResponseWriter, req *http.Request) {
	// http.Redirect(rw, req, "/ui/", 301)
	fmt.Fprintf(rw, "Hello, world - UI")
}

// GetIP gets a requests IP address by reading off the forwarded-for
// header (for proxies) and falls back to use the remote address.
func GetIP(r *http.Request) string {
	forwarded := r.Header.Get("X-FORWARDED-FOR")
	if forwarded != "" {
		return forwarded
	}
	return r.RemoteAddr
}

// ApiVersionMiddleware adds the API version information to the response header
func ApiVersionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		//	Include the version in the response headers:
		rw.Header().Set(version.Header, version.String())

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(rw, r)
	})
}
