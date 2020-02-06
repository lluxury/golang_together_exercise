package main

import (
	log "github.com/Sirupsen/logrus"
	// "github.com/xianlubird/mydocker/container"
	"github.com/lluxury/golang_together_exercise/3_container/container"
	"os"
)

func Run(tty bool, command string) {
	parent := container.NewParentProcess(tty,command)

	if err := parent.Start(); err != nil {
		log.Error(err)
	}

	parent.Wait()
	os.Exit(-1)
}