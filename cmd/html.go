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
	"bufio"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/yosssi/gohtml"
	"os"
	"pretgo/cmd/subcommand"
)

// jsonCmd represents the json command
var htmlCmd = &cobra.Command{
	Use:   "html",
	Short: "This Command is used to pretty json cat response",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			fmt.Println(gohtml.Format(scanner.Text()))
		}
	},
}

func init() {
	htmlCmd.AddCommand(subcommand.FileHtml)
	rootCmd.AddCommand(htmlCmd)
}
