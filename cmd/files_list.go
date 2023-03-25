package cmd

import (
	"fmt"
	"github.com/danesparza/fxaudio/internal/data"
	"github.com/rs/zerolog/log"

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
		log.Err(err).Msg("Problem trying to open the system database")
		return
	}
	defer db.Close()

	//	Get a list of all files:
	gotFiles, err := db.GetAllFiles()
	if err != nil {
		log.Err(err).Msg("Problem trying to get all files")
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
