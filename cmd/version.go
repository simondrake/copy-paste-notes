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
	Short: "Print the version number of copy-paste-notes",
	Run: func(_ *cobra.Command, _ []string) {
		fmt.Println("Copy Paste Notes v0.0.1")
	},
}
