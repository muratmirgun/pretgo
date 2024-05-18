package subcommand

import (
	"log"
	"os"
	"pretgo/internal/json"

	"github.com/spf13/cobra"
)

var FileJson = &cobra.Command{
	Use:   "file [file]",
	Short: "This Command is used to pretty json file",
	Long:  `Example usage: pretgo json file mess.json' `,
	Run: func(cmd *cobra.Command, args []string) {

		if args[0] != "" {
			file, err := os.Open(args[0])
			if err != nil {
				log.Fatal(err)
			}
			defer file.Close()

			if err := json.PrettyFIle(file, os.Stdout); err != nil {
				log.Fatal(err)
			}

		}
	},
}
