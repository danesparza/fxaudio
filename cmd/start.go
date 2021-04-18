package cmd

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/danesparza/fxaudio/api"
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

	//	Create our 'sigs' and 'done' channels (this is for graceful shutdown)
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	//	Indicate what signals we're waiting for:
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)

	//	Create an api service object
	apiService := api.Service{StartTime: time.Now()}

	//	Create a router and setup our REST endpoints...
	UIRouter := mux.NewRouter()

	//	UI ROUTES
	if viper.GetString("uiservice.ui-dir") == "" {
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
		log.Printf("[INFO] Using UI directory: %s\n", viper.GetString("uiservice.ui-dir"))
		UIRouter.PathPrefix("/ui").Handler(http.StripPrefix("/ui", http.FileServer(http.Dir(viper.GetString("uiservice.ui-dir")))))
	}

	/*

		/audio
		PUT - Upload file
		GET - List all files
		POST - Play a random file (or if passed an endpoint in JSON, stream that file)

		/audio/1
		GET - Download file
		POST - Play file
		DELETE - Delete file

	*/
	//	AUDIO ROUTES
	UIRouter.HandleFunc("/v1/audio", apiService.UploadFile).Methods("PUT") //	Upload a file
	UIRouter.HandleFunc("/v1/audio", apiService.PlayAudio).Methods("GET")  // 	List all files
	UIRouter.HandleFunc("/v1/audio", apiService.PlayAudio).Methods("POST") // 	Play a random file (or play the endpoint specified in JSON)

	//	EVENT ROUTES

	//	Setup the CORS options:
	log.Printf("[INFO] Allowed CORS origins: %s\n", viper.GetString("uiservice.allowed-origins"))

	uiCorsRouter := cors.New(cors.Options{
		AllowedOrigins:   strings.Split(viper.GetString("uiservice.allowed-origins"), ","),
		AllowCredentials: true,
	}).Handler(UIRouter)

	//	Format the bound interface:
	formattedUIInterface := viper.GetString("uiservice.bind")
	if formattedUIInterface == "" {
		formattedUIInterface = "127.0.0.1"
	}

	//	Start our shutdown listener (for graceful shutdowns)
	go func() {
		//	If we get a signal...
		_ = <-sigs

		//	Indicate we're done...
		done <- true
	}()

	//	Start the API and UI services
	go func() {
		log.Printf("[INFO] Starting Management UI: http://%s:%s/ui/\n", formattedUIInterface, viper.GetString("uiservice.port"))
		log.Printf("[ERROR] %v\n", http.ListenAndServe(viper.GetString("uiservice.bind")+":"+viper.GetString("uiservice.port"), uiCorsRouter))
	}()

	//	Wait for our signal and shutdown gracefully
	<-done

	log.Printf("[INFO] Shutting down ...")

}

func init() {
	rootCmd.AddCommand(startCmd)
}
