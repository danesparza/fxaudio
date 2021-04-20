package cmd

import (
	"fmt"
	"log"

	"github.com/danesparza/fxaudio/data"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:   "remove [id]",
	Short: "Remove a file managed by fxaudio",
	Args:  cobra.MinimumNArgs(1),
	Long: `
Remove a file managed by fxaudio by passing the ID to remove
	
Example:
fxaudio files remove c1ucgu16v83ji73f1m60`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("remove called with %s\n", args[0])

		//	Create a DBManager object
		db, err := data.NewManager(viper.GetString("datastore.system"))
		if err != nil {
			log.Printf("[ERROR] Error trying to open the system database: %s", err)
			return
		}
		defer db.Close()

		//	Find the file information:
		/*
			gotFiles, err := db.GetAllFiles()
			if err != nil {
				log.Fatalf("[ERROR] Error trying to get all files: %s", err)
			}
		*/

		//	Remove the file from the system

		//	Delete the file on disk

		//	Indicate that the file was removed

		//	Get a list of all files:

	},
}

func init() {
	filesCmd.AddCommand(removeCmd)
}
