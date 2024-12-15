package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Gotainer",
	Long:  `All software has versions. This is Gotainer's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Gotainer simple container runtime v@beta.0.0.1")
	},
}
