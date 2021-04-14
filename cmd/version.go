package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// BuildVersion is the app version number at build time
var BuildVersion = "Unknown"

// CommitID is the SHA commit for the compiled app at build time
var CommitID string

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Shows the version information",
	Run: func(cmd *cobra.Command, args []string) {
		//	Show the version number
		fmt.Printf("\nfxaudio version %s", BuildVersion)

		//	Show the CommitID if available:
		if CommitID != "" {
			fmt.Printf(" (%s)", CommitID[:7])
		}

		//	Trailing space and newline
		fmt.Println(" ")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
