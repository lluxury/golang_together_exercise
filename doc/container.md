[TOC]

### 3 构造容器

####   构造实现run命令版本的容器

#####     Linux proc 文件系统介绍

​      /proc 系统信息
​      相关介绍见 README

#####     实现 run 命令

vi file_map_3.1.md



###### main.go

  定义app
  定义app.Commands 执行2 个命令
    initCommand
    runCommand
  app.Before 定义日志 ？
  app.Run(）



###### main_command.go

  **runCommand**
    Name
    Userage
    Flags 表示指定的参数

​    Action

判断参数，定义变量，Run(tty, cmd)，反回



  **initCommand**  内部函数
    Name
    Userage
    Action

获取参数，使用参数初始化

```go
err := container.RunContainerInitProcess(cmd, nil)
```



###### run.go

  Run() 

定义父进程，调用父进程并等待完成

```go
parent := container.NewParentProcess(tty,command)
if err := parent.Start(); err != nil {
```

由 Start() 方法开始工作：

- 调用创建好的 command
- clone 出来一个 namespace 隔离的进程
- 在子进程中 调用 /proc/self/exe
- 发送 init 参数 ，调用 init 方法 初始化容器 



###### container_process.go

  NewParentProcess() 被调用的父进程函数

```go
func NewParentProcess(tty bool, command string) *exec.Cmd{
	args := []string{"init", command}
	
	cmd := exec.Command("/proc/self/exe", args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags:
		syscall.CLONE_NEWUTS |
		syscall.CLONE_NEWPID |
		syscall.CLONE_NEWNS  |
		syscall.CLONE_NEWNET |
		syscall.CLONE_IPC,
	}
	if tty {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	return cmd
}
```

​      当前进程的工作：

- /proc/self 当前环境，自己调用自己
- init 见初始化函数的 name
- clone 参数，fork 出一个新进程，并用 namesapce 隔离
- 如果指定 -ti 把当前进程的 输入输出导入到标准输入输出 



###### init.go

  RunContainerInitProcess()

```go
func RunContainerInitProcess(command string, args []string) error {
   logrus.Infof("command %s", command)

   defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
   // MS_NOEXEC 本文件系统中不允许运行其他程序
   // MS_NOSUID 运行时 不容许 set-user-ID 或 set-group-ID
   // MS_NODEV Mount 默认设定
    
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
```



#####     实现流程

​      mydocker run
​      解析参数
​      创建隔离容器
​      返回 command 对象?
​      启动进程
​      调用自己
​      初始化, /proc, 执行命令
​      容器开始运行

#####     mydocker

​      container
​        container_process.go
​          NewParentProcess()
​          args, cmd
​          cmd.SysProcAttr 克隆进程
​          tty处理
​        init.go
​          RunContainerInitProcess()
​          定义变量 defaultMountFlags
​          syscall.Mount() 不明
​          syscall.Exec()
​      main_command.go
​        定义变量运行参数 runCommand
​        赋值name,Flags和Action字段,Run(tty, cmd)
​        判断参数,获取命令,调用run 函数
​        定义变量初始化参数 initCommand
​        赋值name,Action字段,使用cli.Context
​        获取传递来的command参数? 容器初始化
​      main.go
​        定义app变量,cli.NewApp() 函数
​        app是结构体,替换Name,Usage属性
​        定义2个命令,使用cli.Command结构体
​        命令之前初始化日志 *cli.Context
​      run.go
​        Run()
​        定义parent, 运行,等待,退出
​    测试
​      go build .
​      ./mydocker run -ti /bin/sh

####   增加容器资源限制

#####     流程

​      创建资源限制容器
​      创建 Subsystem 实例
​      在其hierarchy上配置 cgroup
​      创建 cgroup
​      容器进程进入 cgroup
​      完成资源限制

#####     定义Cgroups的数据结构

​    在启动容器时增加资源限制的配置
​    测试

####   增加管道及环境变量识别

​    管道识别

######       管道

​        进程间通信
​        IPC 的一种
​        半双工，一端写一端读
​        无名管道
​        FIFO管道
​        管道有4KB缓存
​    环境变量识别

​    代码

######       原因

​        特殊字符不方便传递参数
​        有长度限制
​        使用匿名管道

######       流程

​        输入容器内名字
​        传递参数
​        返回 writePipe
​        启动容器进程
​        向writePipe 写与运行命令
​        写入成功
​        执行命令

####   小结

​    创建容器
​    创建父进程
​    挂载文件系统
​      替换init进程
​    完成容器创建