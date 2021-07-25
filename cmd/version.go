package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Jq HTML Version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Jq HTML Version v1.0")
	},
}
