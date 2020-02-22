package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/lluxury/golang_together_exercise/4_image/cgroups/subsystems"
	"strings"

	// "github.com/xianlubird/mydocker/container"
	"github.com/lluxury/golang_together_exercise/4_image/container"
	"os"
)

//func Run(tty bool, command string) {

func Run(tty bool, comArray []string, res *subsystems.ResourceConfig) {

	//parent := container.NewParentProcess(tty,command)
	parent, writePipe := container.NewParentProcess(tty)
	if parent == nil {
		log.Errorf("New parent precess error")
		return
	}

	if err := parent.Start(); err != nil {
		log.Error(err)
	}

	//cgroupManager := cgroups.NewCgroupManager("mydocker-cgroup")
	//defer cgroupManager.Destroy()
	//cgroupManager.Set(res)
	//cgroupManager.Apply(parent.Process.Pid)

	sendInitCommand(comArray, writePipe)
	parent.Wait()

	mntURL := "/root/mnt/"
	rootURL := "/root/"
	container.DeleteWorkSpace(rootURL,mntURL)
	os.Exit(0)
}

func sendInitCommand(comArray []string, writePipe *os.File) {
	command := strings.Join(comArray, " ")
	log.Infof("command all is %s", command)
	writePipe.WriteString(command)
	writePipe.Close()
}
