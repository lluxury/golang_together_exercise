package main

import (
	log "github.com/Sirupsen/logrus"
	"strings"

	// "github.com/xianlubird/mydocker/container"
	"github.com/lluxury/golang_together_exercise/4_image/container"
	"os"
)

//func Run(tty bool, command string) {

//func Run(tty bool, comArray []string, res *subsystems.ResourceConfig) {
func Run(tty bool, comArray []string, volume string) {

	//parent, writePipe := container.NewParentProcess(tty)
	parent, writePipe := container.NewParentProcess(tty,volume)
	if parent == nil {
		log.Errorf("New parent precess error")
		return
	}

	if err := parent.Start(); err != nil {
		log.Error(err)
	}


	sendInitCommand(comArray, writePipe)
	parent.Wait()

	mntURL := "/root/mnt/"
	rootURL := "/root/"
	//container.DeleteWorkSpace(rootURL,mntURL)
	container.DeleteWorkSpace(rootURL,mntURL,volume)
	os.Exit(0)
}

func sendInitCommand(comArray []string, writePipe *os.File) {
	command := strings.Join(comArray, " ")
	log.Infof("command all is %s", command)
	writePipe.WriteString(command)
	writePipe.Close()
}
