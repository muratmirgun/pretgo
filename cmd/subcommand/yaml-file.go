package subcommand

import (
	"github.com/spf13/cobra"
	"log"
	"os"
	"pretgo/internal/yaml"
)

var FileYaml = &cobra.Command{
	Use:   "file [file]",
	Short: "This Command is used to pretty json file",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,

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
