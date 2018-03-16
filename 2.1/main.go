package main

import(
"log"
"os"
"os/exec"
"syscall"
)

func main() {
    cmd :=exec.Command("sh")
    cmd.SysProcAttr = &syscall.SysProcAttr{
        Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWIPC | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWUSER | syscall.CLONE_NEWNET,
    }
    // cmd.SysProcAttr.Credential = &syscall.Credential{Uid: uint32(1), Gid: uint32(1)}

    cmd.SysProcAttr.UidMappings = []syscall.SysProcIDMap{
    {ContainerID: 5001, HostID: syscall.Getuid(), Size: 1},
}
cmd.SysProcAttr.GidMappings = []syscall.SysProcIDMap{
    {ContainerID: 5001, HostID: syscall.Getgid(), Size: 1},
}
    cmd.Stdin = os.Stdin
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr

    if err :=cmd.Run(); err !=nil{
        log.Fatal(err)
    }
    os.Exit(-1)
}



