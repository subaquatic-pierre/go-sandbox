package main

import (
	"log"
	"os"

	"github.com/subaquatic-pierre/gotainer/models"
)

func HandleContainerRun() {
	args := os.Args
	log.Println(args)
	imageName := os.Args[2]

	container, err := models.NewContainer(imageName)
	if err != nil {
		log.Println(err)
		return
	}

	log.Printf("starting container : %s ...\n", container.ID)

	// get command to be passed to container if arg exists
	if len(args) > 2 {
		containerCmd := ""
		containerCmd = args[3]
		container.Cmd = containerCmd
	}

	// mount image to container dir
	container.Init()

	// run container
	container.Exec()

	// handle shutdown on command exit
	container.Shutdown()
}
