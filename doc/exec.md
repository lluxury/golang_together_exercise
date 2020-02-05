调用 linux 名字空间

```go
func main() {
	cmd := exec.Command("sh") // 增加新进程
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWIPC | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWUSER | syscall.CLONE_NEWNET,
	}
	//cmd.SysProcAttr.Credential = &syscall.Credential{Uid: uint32(1), Gid: uint32(1)}
	// CentOS7不支持CLONE_NEWUSER这个flag
    
    cmd.Env = []string{"PS1=-[namespace-process]-# "}  // 环境变量
      // Stdin io.Reader  
	cmd.Stdin = os.Stdin	// 程序的io 管道重定向到标准输入输出
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

    // 运行命令
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
	os.Exit(-1)
}

```

exec 是个包，Command 是exec中函数，传入名字和参数，返回 结构体 Cmd 的指针

 创建 *Cmd 对象，并调用了系统sh命令

Cmd 是个结构体，包含了很多方法

SysProcAttr 参数 是 Cmd 的成员之一，涉及 os.ProcAttr

