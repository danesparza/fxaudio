/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
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
	Run: func(cmd *cobra.Command, args []string) {
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
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
