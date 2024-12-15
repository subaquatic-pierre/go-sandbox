package main

import (
	"os"

	"github.com/subaquatic-pierre/gotainer/cmd"
)

func main() {
	if os.Args[1] == "__run" {
		// main entry point to run container
		HandleContainerRun()
	} else {
		cmd.Execute()
	}
}
