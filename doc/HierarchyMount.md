```go
func main() {
	if os.Args[0] == "/proc/self/exe" {...}
    // 判断 os.Args[0]? 
    // 匹配则 执行命令，信息输出
	// cmd := exec.Command("sh", "-c", `stress --vm-bytes 200m --vm-keep -m 1`)
    
	cmd := exec.Command("/proc/self/exe")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	// 命名空间
    
	if err := cmd.Start(); err != nil {
		fmt.Println("ERROR", err)
		os.Exit(1)
	} else {
		// 得到fork出来进程映射在外部命名空间的pid ?
		fmt.Printf("%v", cmd.Process.Pid)

		os.Mkdir(path.Join(cgroupMemoryHierarchyMount, "testmemorylimit"), 0755)
		ioutil.WriteFile(path.Join(cgroupMemoryHierarchyMount, "testmemorylimit", "tasks"), []byte(strconv.Itoa(cmd.Process.Pid)), 0644)
		ioutil.WriteFile(path.Join(cgroupMemoryHierarchyMount, "testmemorylimit", "memory.limit_in_bytes"), []byte("100m"), 0644)
	}
    // 建目录,写PID, 写限制
	cmd.Process.Wait()
}
```

// 在系统默认创建挂载了memory subsystem的Hierarchy上创建cgroup 

// 将容器进程加入到这个cgroup中

// 限制cgroup进程使用



```go
cmd.Process.Pid
// func Join(elem ...string) string {
// func Mkdir(name string, perm FileMode) error {}
// func WriteFile(filename string, data []byte, perm os.FileMode) error {}
cmd.Process.Wait()

```

