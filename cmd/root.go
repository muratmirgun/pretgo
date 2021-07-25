package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var message string = `
Welcome pretgo
USAGE:
    htmljq [FLAGS] [OPTIONS] <selector>...

FLAGS:
    help                 Prints help information
    phtml                HTML Pretty-print the serialised output attribute <attribute>  Only return this attribute (if present) from selected elements filename <FILE> The input file. Defaults to stdin output <FILE>            The output file. Defaults to stdout
    pxml                 XML Pretty-print the serialised outputs
    Jxml                 JSON Pretty-print the serialised outputs
    version              Prints version information
`
var rootCmd = &cobra.Command{
	Use: "pretgo",
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
