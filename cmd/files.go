package cmd

import (
	"github.com/spf13/cobra"
)

// filesCmd represents the files command
var filesCmd = &cobra.Command{
	Use:   "files",
	Short: "File operations for fxaudio",
	Long: `File operations for fxaudio.  
	
When called without arguments, this will list files managed by fxaudio.`,
	Run: listFiles,
}

func init() {
	rootCmd.AddCommand(filesCmd)
}
