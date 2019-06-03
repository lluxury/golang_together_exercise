package container

import (
	"os"
	"os/exec"
	"syscall"
)

func NewParentProcess(tty bool, command string) *exec.Cmd {
	// 父进程,当前进程执行的内容

	args := []string{"init", command}
	// args 参数, init 是第一个参数, 调用 initCommand
	cmd := exec.Command("/proc/self/exe", args...)
	// /proc/self 运行进程自己的环境, exec.xx 自己调用了自己 对创建进程初始化
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWNET | syscall.CLONE_IPC,
	}
	// fork 出来新进程
	if tty {
		// 判断 ti, 输出 标准输入输出
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	return cmd
}
