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

mydocker

######     测试

```bash
go build .
./3_container run -ti /bin/sh 
ps -ef    
./3_container run -ti /bin/ls
```



####   增加容器资源限制

file_map_3.2

##### 定义Cgroups的数据结构



###### subsystem.go

  tys ResourceConfig  

传递资源限制的结构，内存，cpu时间片，cpu核心



tyi Subsystem Subsystem

接口，定义4个处理方法，cgroup 抽象成 path



  var SubsystemsIns

初始化实例创建资源限制处理链数组

三个数组都还没有



###### memory.go

定义内容实现方法

  tyi MemorySubSystem
  Set
  Remove
  Apply
  Name



###### cpu.go 



###### cpuset.go



```go
func (s *MemorySubSystem) Set(cgroupPath string, res *ResourceConfig) error {
	if subsysCgroupPath, err := GetCgroupPath(s.Name(), cgroupPath, true); err == nil {
		if res.MemoryLimit !="" {
			if err := ioutil.WriteFile(path.Join(subsysCgroupPath, "memory.limit_in_bytes"),[]byte(res.MemoryLimit),0644); err != nil {
				return fmt.Errorf("set cgroup memory fail %v",err)
			}
		}
		return nil
	} else {
		return err
	}
}
```

判断路径，判断内容，判错设置，错误返回

###### utils.go

  **FindCgroupMountpoint( )**

打开文件，搜索内容

```bash
cat /proc/self/mountinfo
# rw,memory ，挂载的 subsystem 是memory
# 在 memory 中增加限制可以限制内存
```

搜索代码块

```go
	scanner := bufio.NewScanner(f)
	for scanner.Scan(){
		txt := scanner.Text()
		fields := strings.Split(txt, " ")
		for _, opt := range strings.Split(fields[len(fields)-1], ",") {
			if opt == subsystem {
				return fields[4]
                //返回 /sys/fs/cgroup/memory
			}
		}
	}
```

  定义bufio搜索，定义Text属性，定义字段，寻找匹配，返回另一部分 



##### 在启动容器时增加资源限制的配置

**GetCgroupPath( )**

返回 Cgroup 的绝对地址

获取挂载点，获取状态，没报错也没有不存在，则返回回值

如果不存在就创建并返回，创建失败就报错

错过获取状态失败就报错



注意体会三层判断，和三级变量指定



###### cgroup_manager.go

管理 Cgroup 与容器建立关系

  tys CgroupManager
  NewCgroupManager( )
  Apply
  Set
  Destroy



###### run.go 新增

  Run()
    tty
    comArray
    res
  cgroupManager.Set()
  cgroupManager.Apply()
  sendInitCommand()
    comArray
    writePipe

下节详细记录



###### main_command.go 

Action 有了较大的变化

```go
	Action: func(context *cli.Context) error {
		if len(context.Args()) < 1 {
			return fmt.Errorf("Missing container command")
		}
		var cmdArray []string
		for _, arg := range context.Args() {
			cmdArray = append(cmdArray, arg)
		}
		tty := context.Bool("ti")
		resConf := &subsystems.ResourceConfig{
		   //MemoryLimit: "m", 重要 bug，直接导致功能失效
			MemoryLimit: context.String("m"),
		}

		Run(tty, cmdArray, resConf)
		return nil
	},
```

Flags  标记增加

initCommand 	命令的调用改变



#####     流程

​      创建资源限制容器
​      创建 Subsystem 实例
​      在其hierarchy上配置 cgroup
​      创建 cgroup
​      容器进程进入 cgroup
​      完成资源限制



##### 测试

测试可以发现一些错误，还是有局限性，只能在环境上使用

代码测试

###### memory_test.go

TestMemoryCgroup( )

memSubSys{}
resConfig {MemoryLimit}
testCgroup 

定义好结构体变量，Cgroup 变量

测试以下方法和函数

Set

FindCgroupMountpoint

Apply

Remove

```go
if err := memSubSys.Set(testCgroup, &resConfig); err != nil {
   t.Fatalf("cgroup fail %v", err)
}
stat, _ := os.Stat(path.Join(FindCgroupMountpoint("memory"),testCgroup))
t.Logf("cgroup stats: %+v", stat)
```



###### utils_test.go

TestFindCgroupMountpoint( )

```go
t.Logf("cpu subsystem mount point %v\n",FindCgroupMountpoint("cpu"))
```

 直接调用 FindCgroupMountpoint 把挂载点打印出来



###### 命令测试

```bash
cd 3_container/cgroups/subsystems
go test -v

./3_container run -ti -m 100m stress \
	--vm-bytes 200m --vm-keep -m 1

./3_container run -ti -cpushare 512 stress --vm-bytes 200m --vm-keep -m 1
```





####   增加管道及环境变量识别

**管道**

​        进程间通信
​        IPC 的一种
​        半双工，一端写一端读
​        无名管道
​        FIFO管道
​        管道有4KB缓存
​    环境变量识别



**原因**

​        特殊字符不方便传递参数
​        有长度限制
​        使用匿名管道



​    **代码**

###### container_process.go 改

NewParentProcess()  修改了参数个数，新增了功能

NewPipe()	pipe 方法生成匿名管道

```go
readPipe, writePipe , err := NewPipe()
cmd.ExtraFiles = []*os.File{readPipe}
// 外带文件句柄创建子进程,传入管道文件的读取端句柄

func NewPipe()(*os.File, *os.File, error){
	read, write,err := os.Pipe()
	if err != nil {
		return nil, nil, err
	}
	return read, write, nil
}
```



补充

```bash
ll /proc/self/fd
lrwx------ 1 root root 64 Feb  8 20:05 0 -> /dev/pts/0
lrwx------ 1 root root 64 Feb  8 20:05 1 -> /dev/pts/0
lrwx------ 1 root root 64 Feb  8 20:05 2 -> /dev/pts/0
lr-x------ 1 root root 64 Feb  8 20:05 3 -> /proc/9196/fd/
                                         # 传给了进程
```



###### container/init.go

RunContainerInitProcess()   逻辑改变，函数重写

```go
cmdArray := readUserCommand()
exec.LookPath(cmdArray[0])   //识别环境变量

syscall.Exec(path, cmdArray[0:], os.Environ()); // 调用自己
```



readUserCommand()	新增 read 函数

```go
func readUserCommand() []string {
	pipe := os.NewFile(uintptr(3), "pipe") // 第四描述符
	msg, err := ioutil.ReadAll(pipe) // 自己维护等待
	if err != nil {
		logrus.Errorf("init read pipe error %v", err)
		return nil
	}
	msgStr := string(msg)
	return strings.Split(msgStr, " ")
}
// 获取并分割

```



###### run.go改

run 函数 参数个数都变了：

- ​    tty
- ​    comArry
- ​    res

  进程启动完再发送命令



  **sendInitCommand( )**
    发送命令的函数

```go
parent, writePipe := container.NewParentProcess(tty) //?
sendInitCommand(comArry, writePipe)

func sendInitCommand(comArray []string, writePipe *os.File) {
	command := strings.Join(comArray, " ")
	log.Infof("command all is %s", command)
	writePipe.WriteString(command)
	writePipe.Close()
}
```



###### 测试

```bash
./3_container  run -ti ls -l
./3_container  run -ti bash
 ps -ef
```



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