package subcommand

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/yosssi/gohtml"
)

var FileHtml = &cobra.Command{
	Use:   "file [file]",
	Short: "This Command is used to pretty html file",
	Long:  `Example usage: pretgo html file mess.html' `,
	Run: func(cmd *cobra.Command, args []string) {
		if args[0] != "" {
			file, err := os.Open(args[0])
			if err != nil {
				log.Fatal(err)
			}
			defer file.Close()

			scanner := bufio.NewScanner(os.Stdin)
			for scanner.Scan() {
				fmt.Println(gohtml.Format(scanner.Text()))
			}

		}

	},
}
