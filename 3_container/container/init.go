package container

import (
	"os"
	"syscall"

	"github.com/Sirupsen/logrus"
)

func RunContainerInitProcess(command string, args []string) error {
	logrus.Infof("command %s", command)

	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	// MS_NOEXEC 本文件系统中不允许运行其他程序
	// MS_NOSUID 运行时 不容许 set-user-ID 或 set-group-ID
	// MS_NODEV 默认设定
	syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")
	argv := []string{command}
	if err := syscall.Exec(command, argv, os.Environ()); err != nil {
		// 容器第一个程序,是init初始化进程,不能kill,
		// 调用 int execve(const char *filename, char *const argv[],char *constenvp[];)
		// 执行当前 filename对应的程序, 覆盖当前镜像,数据和堆栈,
		logrus.Errorf(err.Error())
	}
	return nil
}
