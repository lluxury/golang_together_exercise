package container

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/Sirupsen/logrus"
)


//func RunContainerInitProcess(command string, args []string) error {
func RunContainerInitProcess() error {
	//logrus.Infof("command %s", command)
	cmdArray := readUserCommand()
	if cmdArray == nil || len(cmdArray) == 0 {
		return fmt.Errorf("Run container get user command error, cmdArray is nil")
	}

	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")

	//setUpMount()
	path, err := exec.LookPath(cmdArray[0])
	if err != nil {
		logrus.Errorf("Exec loop path error %v", err)
		return err
	}
	logrus.Info("Find path %s", path)


	//argv :=[]string{command}

	//if err := syscall.Exec(command, argv, os.Environ()); err != nil{
	//	logrus.Errorf(err.Error())
	//}

	if err := syscall.Exec(path, cmdArray[0:], os.Environ()); err != nil {
		logrus.Errorf(err.Error())
	}

	return nil
}

func readUserCommand() []string {
	pipe := os.NewFile(uintptr(3), "pipe")
	msg, err := ioutil.ReadAll(pipe)
	if err != nil {
		logrus.Errorf("init read pipe error %v", err)
		return nil
	}
	msgStr := string(msg)
	return strings.Split(msgStr, " ")
}


