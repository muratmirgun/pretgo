package subcommand

import (
	"log"
	"os"
	"pretgo/internal/xml"

	"github.com/spf13/cobra"
)

var FileXml = &cobra.Command{
	Use:   "file [file]",
	Short: "This Command is used to pretty xml file",
	Long:  `Example usage: pretgo json file mess.xml' `,

	Run: func(cmd *cobra.Command, args []string) {

		if args[0] != "" {
			file, err := os.Open(args[0])
			if err != nil {
				log.Fatal(err)
			}
			defer file.Close()

			if err := xml.Pretty(file, os.Stdout); err != nil {
				log.Fatal(err)
			}
		}
	},
}
