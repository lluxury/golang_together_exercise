[toc]

### 用户 id 相关问题

#### CLONE_NEWUSER

问题，添加该标记后，丢失 root 用户



#### UID和GID的映射

- 由于 USER 命名空间，UID 可以和环境分离
- 主机上有多个 USER 命名空间使用。
- 一个进程可以在**多个**命名空间中使用**不同**的 UID
- UID和GID映射解决不同命名空间下的ID匹配问题 



#### 操作映射

```go
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags:
		syscall.CLONE_NEWNS  |
			syscall.CLONE_NEWUTS |
			syscall.CLONE_NEWIPC |
			syscall.CLONE_NEWPID |
			syscall.CLONE_NEWNET |
			syscall.CLONE_NEWUSER,
			UidMappings: []syscall.SysProcIDMap{
			{
				ContainerID :0,
				HostID: os.Getuid(),
				Size:  1,
			},
			},

		GidMappings: []syscall.SysProcIDMap{
			{
				ContainerID :0,
				HostID: os.Getgid(),
				Size:  1,
			},
		},
	}
```

以上为调用片段，完整代码看 github



cmd.SysProcAttr 有两个属性 **UidMappings** 和 **GidMappings**，定义类型

```go
type SysProcIDMap struct {
    ContainerID int // Container ID. 容器内
    HostID      int // Host ID. 宿主机内
    Size        int // Size.  范围
}
```

