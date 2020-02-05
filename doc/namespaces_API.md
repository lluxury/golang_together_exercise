[toc]

### namespaces API

以下3个都是 **api**

- clone()   创建新进程
- unshare()   进程移出 Namespace
- setns()   进程加入 Namespace



#### unshare 命令

```bash
unshare -h
unshare -u /bin/sh  # 进入了新的名字空间 --uts
 # 类似 syscall.CLONE_NEWUTS
```



#### clone( )

这个api还能接受其他的参数:

```bash
    CLONE_IO
    CLONE_NEWIPC
    CLONE_NEWNET
    CLONE_NEWNS 
    ...
```

运行程序，转化成 调用Linux namespace的 clone( ) api，

并传递了  CLONE_* 参数来创建新进程， 例如 /bin/sh



#### 调试

```bash
strace /tmp/main 
execve("/tmp/main", ["/tmp/main"], [/* 32 vars */]) = 0
arch_prctl(ARCH_SET_FS, 0x58e490)       = 0
... mmap
clone(child_stack=0, flags=CLONE_NEWNS|CLONE_NEWUTS|CLONE_NEWIPC|CLONE_NEWUSER|CLONE_NEWPID|CLONE_NEWNET|SIGCHLD) = 3046
...
waitid(P_PID, 3046, $  

# 未结束。。

readlink /proc/self/ns/uts  # 测试
```

  