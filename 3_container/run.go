package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/lluxury/golang_together_exercise/3_container/cgroups/subsystems"
	"strings"

	// "github.com/xianlubird/mydocker/container"
	"github.com/lluxury/golang_together_exercise/3_container/container"
	"os"
)

//func Run(tty bool, command string) {
func Run(tty bool, comArry []string, res *subsystems.ResourceConfig)  {

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

	sendInitCommand(comArry, writePipe)
	parent.Wait()
	os.Exit(0)
}

func sendInitCommand(comArray []string, writePipe *os.File) {
	command := strings.Join(comArray, " ")
	log.Infof("command all is %s", command)
	writePipe.WriteString(command)
	writePipe.Close()
}