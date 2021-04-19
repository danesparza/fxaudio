package cmd

import (
	"fmt"
	"log"

	"github.com/danesparza/fxaudio/data"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List files managed by fxaudio",
	Long:  `List files managed by fxaudio`,
	Run:   listFiles,
}

func listFiles(cmd *cobra.Command, args []string) {
	//	Create a DBManager object
	db, err := data.NewManager(viper.GetString("datastore.system"))
	if err != nil {
		log.Printf("[ERROR] Error trying to open the system database: %s", err)
		return
	}
	defer db.Close()

	//	Get a list of all files:
	gotFiles, err := db.GetAllFiles()
	if err != nil {
		log.Fatalf("[ERROR] Error trying to get all files: %s", err)
	}

	//	Gather our formatted output
	formattedList := pterm.TableData{{"ID", "FilePath"}}
	for _, v := range gotFiles {
		formattedList = append(formattedList, []string{v.ID, v.FilePath})
	}

	//	Render the output:
	fmt.Println()
	pterm.DefaultTable.WithHasHeader().WithData(formattedList).Render()
	fmt.Println()
}

func init() {
	filesCmd.AddCommand(listCmd)

}
