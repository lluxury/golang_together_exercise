package main

import(
"log"
"os"
"os/exec"
"syscall"
)

func main() {
    cmd :=exec.Command("sh")    // 新进程初始命令
    cmd.SysProcAttr = &syscall.SysProcAttr{
        Cloneflags: syscall.CLONE_NEWUTS,
        // syscall.CLONE_NEWUTS 创建 UTS Namespace等

    }
    //}
    cmd.Stdin = os.Stdin
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr

    if err :=cmd.Run(); err !=nil{
        log.Fatal(err)
    }
    // os.Exit(-1)
}



