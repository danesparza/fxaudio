package cmd

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/hashicorp/logutils"
	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var (
	cfgFile  string
	loglevel string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "fxaudio",
	Short: "REST service for multichannel audio",
	Long:  `REST service for multichannel audio.  For your digital effects and soundtrack needs.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/fxaudio.yaml)")
	rootCmd.PersistentFlags().StringVarP(&loglevel, "loglevel", "l", "WARN", "Log level: DEBUG/INFO/WARN/ERROR")

	//	Bind config flags for optional config file override:
	viper.BindPFlag("loglevel", rootCmd.PersistentFlags().Lookup("loglevel"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// Find home directory.
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(home)      // adding home directory as first search path
		viper.AddConfigPath(".")       // also look in the working directory
		viper.SetConfigName("fxaudio") // name the config file (without extension)
	}

	viper.AutomaticEnv() // read in environment variables that match

	//	Set our defaults
	viper.SetDefault("loglevel", "ERROR")
	viper.SetDefault("datastore.system", path.Join(home, "fxaudio", "db", "system.db"))
	viper.SetDefault("datastore.retentiondays", 30)
	viper.SetDefault("uiservice.port", 3000)
	viper.SetDefault("uiservice.allowed-origins", "*")

	// Use the config file.  Track if there was an error
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("[ERROR] There was an error reading in the config file: %v", err)
	}

	//	Set the log level from config (if we have it)
	filter := &logutils.LevelFilter{
		Levels:   []logutils.LogLevel{"DEBUG", "INFO", "WARN", "ERROR"},
		MinLevel: logutils.LogLevel(viper.GetString("loglevel")),
		Writer:   os.Stderr,
	}
	log.SetOutput(filter)
}
