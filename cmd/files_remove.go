package cmd

import (
	"fmt"
	"github.com/danesparza/fxaudio/internal/data"
	"github.com/rs/zerolog/log"
	"os"

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
		removeId := args[0]

		//	Create a DBManager object
		db, err := data.NewManager(viper.GetString("datastore.system"))
		if err != nil {
			log.Err(err).Msg("Problem trying to open the system database")
			return
		}
		defer db.Close()

		//	Find the file information:
		gotFile, err := db.GetFile(removeId)
		if err != nil {
			log.Err(err).Str("fileid", removeId).Msg("Problem trying to find fileid")
			return
		}

		//	Remove the file from the system
		if err = db.DeleteFile(removeId); err != nil {
			log.Err(err).Str("file", gotFile.FilePath).Msg("Problem removing file from system")
			return
		}

		//	Delete the file on disk
		if err = os.Remove(gotFile.FilePath); err != nil {
			log.Err(err).Str("file", gotFile.FilePath).Msg("Problem removing file from disk")
			return
		}

		//	Indicate that the file was removed
		fmt.Printf("\nFile %s removed\n", removeId)

	},
}

func init() {
	filesCmd.AddCommand(removeCmd)
}
