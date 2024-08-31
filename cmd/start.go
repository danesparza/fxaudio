package cmd

import (
	"context"
	"fmt"
	"github.com/danesparza/fxaudio/internal/data"
	"github.com/danesparza/fxaudio/internal/media"
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

	//	Trap program exit appropriately
	ctx, cancel := context.WithCancel(context.Background())
	sigs := make(chan os.Signal, 2)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)
	go handleSignals(ctx, sigs, cancel)

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

	//	Get the system config:
	systemConfig, err := appdata.GetConfig(ctx)
	if err != nil {
		log.Err(err).Msg("Problem trying to read the system config")
		return
	}

	audioSvc := media.NewVLCAudioService(systemConfig.AlsaDevice)

	//	Create a background service object
	backgroundService := media.BackgroundProcess{
		PlayAudio:    make(chan media.PlayAudioRequest),
		StopAudio:    make(chan string),
		StopAllAudio: make(chan bool),
		DB:           appdata,
		AS:           audioSvc,
	}

	//	Create an api service object
	apiService := api.Service{
		PlayMedia:    backgroundService.PlayAudio,
		StopMedia:    backgroundService.StopAudio,
		StopAllMedia: backgroundService.StopAllAudio,
		DB:           appdata,
		AS:           audioSvc,
		StartTime:    time.Now(),
	}

	//	Set up the API routes
	r := api.NewRouter(apiService)

	//	Start the media processor:
	go backgroundService.HandleAndProcess(ctx)

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
