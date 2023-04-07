package cmd

import (
	"context"
	"fmt"
	"github.com/danesparza/fxaudio/internal/data"
	"github.com/danesparza/fxaudio/internal/media"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/danesparza/fxaudio/api"
	_ "github.com/danesparza/fxaudio/docs" // swagger docs location
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	httpSwagger "github.com/swaggo/http-swagger" // http-swagger middleware
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the API and UI services",
	Long:  `Start the API and UI services`,
	Run:   start,
}

func start(cmd *cobra.Command, args []string) {

	//	If we have a config file, report it:
	if viper.ConfigFileUsed() != "" {
		log.Debug().Str("configFile", viper.ConfigFileUsed()).Msg("Using config file")
	} else {
		log.Debug().Msg("No config file found")
	}

	systemdb := viper.GetString("datastore.system")
	uploadPath := viper.GetString("upload.path")
	uploadByteLimit := viper.GetString("upload.bytelimit")

	//	Emit what we know:
	log.Info().
		Str("systemdb", systemdb).
		Str("uploadPath", uploadPath).
		Str("uploadByteLimit", uploadByteLimit).
		Msg("Config")

	//	Init SQLite
	db, err := data.InitSqlite(systemdb)
	if err != nil {
		log.Err(err).Msg("Problem trying to open the system database")
		return
	}
	defer db.Close()

	//	Init the AppDataService
	appdata := data.NewAppDataService(db)

	//	Create an api service object
	apiService := api.Service{
		PlayMedia:    make(chan media.PlayAudioRequest),
		StopMedia:    make(chan string),
		StopAllMedia: make(chan bool),
		DB:           appdata,
		StartTime:    time.Now(),
	}

	//	Trap program exit appropriately
	ctx, cancel := context.WithCancel(context.Background())
	sigs := make(chan os.Signal, 2)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)
	go handleSignals(ctx, sigs, cancel)

	//	Create a router and set up our REST endpoints...
	r := chi.NewRouter()

	//	Add middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Compress(5))
	r.Use(api.ApiVersionMiddleware)

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
		r.Post("/audio/play/random/{tag}", apiService.PlayRandomAudio)
		r.Post("/audio/stream", apiService.StreamAudio)
		r.Post("/audio/loop/{id}/{loopTimes}", apiService.LoopAudio)

		//	Stop audio
		r.Post("/audio/stop", apiService.StopAllAudio)
		r.Post("/audio/stop/{id}", apiService.StopAudio)

	})

	//	SWAGGER
	r.Mount("/swagger", httpSwagger.WrapHandler)

	//	Start the media processor:
	go media.HandleAndProcess(ctx, apiService.PlayMedia, apiService.StopMedia, apiService.StopAllMedia)

	//	Format the bound interface:
	formattedServerPort := fmt.Sprintf(":%v", viper.GetString("server.port"))

	//	Start the service and display how to access it
	log.Info().Str("server", formattedServerPort).Msg("Started REST service")
	log.Err(http.ListenAndServe(formattedServerPort, r)).Msg("HTTP API service error")
}

func handleSignals(ctx context.Context, sigs <-chan os.Signal, cancel context.CancelFunc) {
	select {
	case <-ctx.Done():
	case sig := <-sigs:
		switch sig {
		case os.Interrupt:
			log.Info().Msg("SIGINT")
		case syscall.SIGTERM:
			log.Info().Msg("SIGTERM")
		}

		log.Info().Msg("Shutting down ...")
		cancel()
		os.Exit(0)
	}
}

func init() {
	rootCmd.AddCommand(startCmd)
}
