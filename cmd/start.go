package cmd

import (
	"context"
	"fmt"
	"github.com/danesparza/fxaudio/internal/data"
	"github.com/danesparza/fxaudio/internal/event"
	"github.com/danesparza/fxaudio/internal/media"
	"github.com/danesparza/fxaudio/internal/mediatype"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"net/http"
	"os"
	"os/signal"
	"strconv"
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

	retentiondays := viper.GetString("datastore.retentiondays")
	systemdb := viper.GetString("datastore.system")
	uploadPath := viper.GetString("upload.path")
	uploadByteLimit := viper.GetString("upload.bytelimit")

	//	Emit what we know:
	log.Info().
		Str("systemdb", systemdb).
		Str("retentiondays", retentiondays).
		Str("uploadPath", uploadPath).
		Str("uploadByteLimit", uploadByteLimit).
		Msg("Config")

	//	Parse the log retention (in days):
	historyttl, err := strconv.Atoi(retentiondays)
	if err != nil {
		log.Err(err).Msg("The datastore.retentiondays config is invalid")
		return
	}

	//	Create a DBManager object and associate with the api.Service
	db, err := data.NewManager(systemdb)
	if err != nil {
		log.Err(err).Msg("Problem trying to open the system database")
		return
	}
	defer db.Close()

	//	Create an api service object
	apiService := api.Service{
		PlayMedia:    make(chan media.PlayAudioRequest),
		StopMedia:    make(chan string),
		StopAllMedia: make(chan bool),
		DB:           db,
		StartTime:    time.Now(),
		HistoryTTL:   time.Duration(int(historyttl)*24) * time.Hour,
	}

	//	Trap program exit appropriately
	ctx, cancel := context.WithCancel(context.Background())
	sigs := make(chan os.Signal, 2)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)
	go handleSignals(ctx, sigs, cancel, db, apiService.HistoryTTL)

	//	Log that the system has started:
	_, err = db.AddEvent(event.SystemStartup, mediatype.System, "System started", "", apiService.HistoryTTL)
	if err != nil {
		log.Err(err).Msg("Problem trying to log to the system datastore")
		return
	}

	//	Create a router and setup our REST endpoints...
	r := chi.NewRouter()

	//	Add middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Compress(5))

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
		//	Audio routes
		r.Put("/audio", apiService.UploadFile)
		r.Get("/audio", apiService.ListAllFiles)
		r.Delete("/audio/{id}", apiService.DeleteFile)

		r.Post("/audio/play", apiService.PlayRandomAudio)
		r.Post("/audio/play/{id}", apiService.PlayAudio)
		r.Post("/audio/stream", apiService.StreamAudio)
		r.Post("/audio/stop", apiService.StopAllAudio)
		r.Post("/audio/stop/{id}", apiService.StopAudio)

		//	Event routes
		r.Get("/events", apiService.GetAllEvents)
		r.Get("/event/{id}", apiService.GetEvent)
	})

	//	SWAGGER
	r.Mount("/swagger", httpSwagger.WrapHandler)

	//	Start the media processor:
	go media.HandleAndProcess(ctx, apiService.PlayMedia, apiService.StopMedia, apiService.StopAllMedia)

	//	Format the bound interface:
	formattedServerInterface := viper.GetString("server.bind")
	if formattedServerInterface == "" {
		formattedServerInterface = GetOutboundIP().String()
	}

	//	Start the service and display how to access it
	log.Info().
		Str("documentation-url", fmt.Sprintf("http://%s:%s/swagger/", formattedServerInterface, viper.GetString("server.port"))).
		Msg("REST service started")

	err = http.ListenAndServe(viper.GetString("server.bind")+":"+viper.GetString("server.port"), r)
	if err != nil {
		log.Err(err).Msg("Problem with REST server")
	}
}

func handleSignals(ctx context.Context, sigs <-chan os.Signal, cancel context.CancelFunc, db *data.Manager, historyttl time.Duration) {
	select {
	case <-ctx.Done():
	case sig := <-sigs:
		switch sig {
		case os.Interrupt:
			log.Info().Msg("SIGINT")
		case syscall.SIGTERM:
			log.Info().Msg("SIGTERM")
		}

		//	Log that the system has started:
		_, err := db.AddEvent(event.SystemShutdown, mediatype.System, "System stopping", "", historyttl)
		if err != nil {
			log.Printf("[ERROR] Error trying to log to the system datastore: %s", err)
		}

		log.Info().Msg("Shutting down ...")
		cancel()
		os.Exit(0)
	}
}

func init() {
	rootCmd.AddCommand(startCmd)
}
