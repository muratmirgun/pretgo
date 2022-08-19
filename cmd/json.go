// Package cmd
/*

Copyright © 2022 Murat

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
	"log"
	"os"
	"pretgo/cmd/subcommand"
	"pretgo/internal/json"

	"github.com/spf13/cobra"
)

// jsonCmd represents the json command
var jsonCmd = &cobra.Command{
	Use:   "json",
	Short: "This Command is used to pretty json file cat ",
	Long:  `Example usage: cat mess.json | pretgo json `,
	Run: func(cmd *cobra.Command, args []string) {

		if err := json.Pretty(os.Stdin, os.Stdout); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	jsonCmd.AddCommand(subcommand.FileJson)
	rootCmd.AddCommand(jsonCmd)
}
