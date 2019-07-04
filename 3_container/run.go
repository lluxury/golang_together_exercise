package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/xianlubird/mydocker/container"
	"os"
)

func Run(tty bool, command string) {
	parent := container.NewParentProcess(tty, command)
	if err := parent.Start(); err != nil {
		log.Error(err)
	}
	parent.Wait()
	os.Exit(-1)
	// start 方法是创建好 command 的调用,克隆进程,在子进程中调用自己,发送init参数去初始化资源
}
