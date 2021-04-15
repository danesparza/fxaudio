package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	jsonConfig bool
	yamlConfig bool
)

var yamlDefault = []byte(`
loglevel: INFO
`)

var jsonDefault = []byte(`{	
	"loglevel": "INFO"
}`)

// defaultsCmd represents the defaults command
var defaultsCmd = &cobra.Command{
	Use:   "defaults",
	Short: "Prints default fxaudio configuration file",
	Long: `Use this to create a default configuration file for fxaudio. 
	Example:
	fxaudio defaults > fxaudio.yaml`,
	Run: func(cmd *cobra.Command, args []string) {
		if jsonConfig {
			fmt.Printf("%s", jsonDefault)
		} else if yamlConfig {
			fmt.Printf("%s", yamlDefault)
		}
	},
}

func init() {
	rootCmd.AddCommand(defaultsCmd)

	defaultsCmd.Flags().BoolVarP(&jsonConfig, "json", "j", false, "Create a JSON configuration file")
	defaultsCmd.Flags().BoolVarP(&yamlConfig, "yaml", "y", true, "Create a YAML configuration file")
}
