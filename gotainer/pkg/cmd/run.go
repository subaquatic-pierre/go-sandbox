package cmd

import (
	"os"
	"os/exec"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/subaquatic-pierre/gotainer/util"
)

var runCmd = &cobra.Command{
	Use:   "run [image]",
	Short: "Run a container with the given image name",
	Long: `This command starts a new container with the given image
              name and command if specified. If the container tar doesn't
                exist then the command fails`,
	Example: "run ubuntu /bin/bash",
	Run: func(_cmd *cobra.Command, args []string) {
		// show usage if now args passed
		if len(args) == 0 {
			_cmd.Help()
			return
		}

		// re-run this, essentially forking the process
		cmd := exec.Command("/proc/self/exe", append([]string{"__run"}, os.Args[2:]...)...)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.SysProcAttr = &syscall.SysProcAttr{
			Cloneflags:   syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
			Unshareflags: syscall.CLONE_NEWNS,
		}

		util.Must(cmd.Run())

	},
}
