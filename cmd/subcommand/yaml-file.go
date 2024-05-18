package subcommand

import (
	"log"
	"os"
	"pretgo/internal/yaml"

	"github.com/spf13/cobra"
)

var FileYaml = &cobra.Command{
	Use:   "file [file]",
	Short: "This Command is used to pretty yaml file",
	Long:  `Example usage: pretgo json file mess.yaml' `,

	Run: func(cmd *cobra.Command, args []string) {

		if args[0] != "" {
			file, err := os.Open(args[0])
			if err != nil {
				log.Fatal(err)
			}
			defer file.Close()

			if err := yaml.Pretty(file); err != nil {
				log.Fatal(err)
			}
		}
	},
}
