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
	"pretgo/internal/yaml"

	"github.com/spf13/cobra"
)

// jsonCmd represents the json command
var yamlCmd = &cobra.Command{
	Use:   "yaml",
	Short: "This Command is used to pretty yaml file cat ",
	Long:  `Example usage: cat mess.yaml | pretgo yaml `,
	Run: func(cmd *cobra.Command, args []string) {
		if err := yaml.Pretty(os.Stdin); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	yamlCmd.AddCommand(subcommand.FileYaml)
	rootCmd.AddCommand(yamlCmd)
}
