package cmd

import (
	"github.com/spf13/cobra"
	"github.com/subaquatic-pierre/gotainer/lib"
)

var pullCmd = &cobra.Command{
	Use:   "pull [image]",
	Short: "Pull an image with given name from docker",
	Long: `Pull an image from docker or from any other repository. The image will 
          be stored in ./assets/{image}.tar.gz. Which can then be used as file system
          for running that container`,
	Example: "pull ubuntu",
	Run: func(cmd *cobra.Command, args []string) {
		// show usage if now args passed
		if len(args) == 0 {
			cmd.Help()
			return
		}

		// get image name and command to be used
		imageName := args[0]

		lib.PullImage(imageName)

	},
}
