package cmd

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/danesparza/fxaudio/api"
	"github.com/danesparza/fxaudio/data"
	"github.com/danesparza/fxaudio/media"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
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
		log.Println("[DEBUG] Using config file:", viper.ConfigFileUsed())
	} else {
		log.Println("[DEBUG] No config file found.")
	}

	//	Trap program exit appropriately
	ctx, cancel := context.WithCancel(context.Background())
	sigs := make(chan os.Signal, 2)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)
	go handleSignals(ctx, sigs, cancel)

	//	Create a DBManager object and associate with the api.Service
	db, err := data.NewManager(viper.GetString("datastore.system"))
	if err != nil {
		log.Printf("[ERROR] Error trying to open the system database: %s", err)
		return
	}
	defer db.Close()

	//	Create an api service object
	apiService := api.Service{PlayMedia: make(chan media.PlayAudioRequest), StopMedia: make(chan string), StopAllMedia: make(chan bool), DB: db, StartTime: time.Now()}

	//	Create a router and setup our REST endpoints...
	restRouter := mux.NewRouter()

	//	UI ROUTES
	if viper.GetString("server.ui-dir") == "" {
		//	Use the static assets file generated with
		//	https://github.com/elazarl/go-bindata-assetfs using the application-monitor-ui from
		//	https://github.com/danesparza/application-monitor-ui.
		//
		//	To generate this file, run `yarn build` under the "navajo-plex-ui" project.
		//	Then rename the 'build' directory to 'ui', place that
		//	directory under the main navajo-plex directory and run the commands:
		//	go-bindata-assetfs -pkg cmd -o .\cmd\bindata.go ./ui/...
		//	go install ./...

		//  UIRouter.PathPrefix("/ui").Handler(http.StripPrefix("/ui", http.FileServer(assetFS())))
	} else {
		//	Use the supplied directory:
		log.Printf("[INFO] Using UI directory: %s\n", viper.GetString("server.ui-dir"))
		restRouter.PathPrefix("/ui").Handler(http.StripPrefix("/ui", http.FileServer(http.Dir(viper.GetString("server.ui-dir")))))
	}

	//	AUDIO ROUTES
	restRouter.HandleFunc("/v1/audio", apiService.UploadFile).Methods("PUT")         // Upload a file
	restRouter.HandleFunc("/v1/audio", apiService.ListAllFiles).Methods("GET")       // List all files
	restRouter.HandleFunc("/v1/audio/{id}", apiService.DeleteFile).Methods("DELETE") // Delete a file

	restRouter.HandleFunc("/v1/audio/play", apiService.PlayRandomAudio).Methods("POST") // Play a random file
	restRouter.HandleFunc("/v1/audio/play/{id}", apiService.PlayAudio).Methods("POST")  // Play a file
	restRouter.HandleFunc("/v1/audio/stream", apiService.StreamAudio).Methods("POST")   // Stream the endpoint specified in JSON
	restRouter.HandleFunc("/v1/audio/stop", apiService.StopAllAudio).Methods("POST")    // Stop playing all audio
	restRouter.HandleFunc("/v1/audio/stop/{pid}", apiService.StopAudio).Methods("POST") // Stop playing a file

	//	EVENT ROUTES
	//	/v1/events				//	List all events from the past n days

	//	Start the media processor:
	go media.HandleAndProcess(ctx, apiService.PlayMedia, apiService.StopMedia, apiService.StopAllMedia)

	//	Setup the CORS options:
	log.Printf("[INFO] Allowed CORS origins: %s\n", viper.GetString("server.allowed-origins"))

	uiCorsRouter := cors.New(cors.Options{
		AllowedOrigins:   strings.Split(viper.GetString("server.allowed-origins"), ","),
		AllowCredentials: true,
	}).Handler(restRouter)

	//	Format the bound interface:
	formattedServerInterface := viper.GetString("server.bind")
	if formattedServerInterface == "" {
		formattedServerInterface = "127.0.0.1"
	}

	//	Start the API and UI services
	log.Printf("[INFO] Starting Management UI: http://%s:%s/ui/\n", formattedServerInterface, viper.GetString("server.port"))
	log.Printf("[ERROR] %v\n", http.ListenAndServe(viper.GetString("server.bind")+":"+viper.GetString("server.port"), uiCorsRouter))
}

func handleSignals(ctx context.Context, sigs <-chan os.Signal, cancel context.CancelFunc) {
	select {
	case <-ctx.Done():
	case sig := <-sigs:
		switch sig {
		case os.Interrupt:
			log.Println("[INFO] SIGINT")
		case syscall.SIGTERM:
			log.Println("[INFO] SIGTERM")
		}
		log.Println("[INFO] Shutting down ...")
		cancel()
		os.Exit(0)
	}
}

func init() {
	rootCmd.AddCommand(startCmd)
}
