package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const versionNumber = "v0.3.0"

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Prints the version number",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(versionNumber)
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
