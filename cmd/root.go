package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var message string = `
Welcome jq for HTML
USAGE:
    htmljq [FLAGS] [OPTIONS] <selector>...

FLAGS:
    help                 Prints help information
    pretty               Pretty-print the serialised output
    text                 Output only the contents of text nodes inside selected elements
    version              Prints version information

OPTIONS:
    attribute <attribute>    Only return this attribute (if present) from selected elements
    filename <FILE>          The input file. Defaults to stdin
    output <FILE>            The output file. Defaults to stdout
`
var rootCmd = &cobra.Command{
	Use: "htmljq",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(message)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
