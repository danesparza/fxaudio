package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
)

func NewRouter(apiService Service) http.Handler {
	//	Create a router and set up our REST endpoints...
	r := chi.NewRouter()

	//	Add middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Compress(5))
	r.Use(ApiVersionMiddleware)

	//	... including CORS middleware
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	r.Route("/v1", func(r chi.Router) {
		//	File management
		r.Put("/audio", apiService.UploadFile)
		r.Get("/audio", apiService.ListAllFiles)
		r.Delete("/audio/{id}", apiService.DeleteFile)
		r.Post("/audio/{id}", apiService.UpdateTags)

		//	Play audio
		r.Post("/audio/play", apiService.PlayRandomAudio)
		r.Post("/audio/play/{id}", apiService.PlayAudio)
		r.Post("/audio/play/random/{tag}", apiService.PlayRandomAudioWithTag)
		r.Post("/audio/stream", apiService.StreamAudio)
		r.Post("/audio/loop/{id}/{loopTimes}", apiService.LoopAudio)

		//	Stop audio
		r.Post("/audio/stop", apiService.StopAllAudio)
		r.Post("/audio/stop/{id}", apiService.StopAudio)
	})

	//	SWAGGER
	r.Mount("/swagger", httpSwagger.WrapHandler)

	return r
}
