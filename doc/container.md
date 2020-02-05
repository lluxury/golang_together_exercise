[TOC]

### 3 构造容器

####   构造实现run命令版本的容器

#####     Linux proc 文件系统介绍

​      /proc 系统信息
​      相关介绍见 README

#####     实现 run 命令

​      父进程
​        /proc/self/exe
​          进程环境
​          ,exe 进程自己
​        初始化进程
​        clone 参数, fork 新进程,隔离
​        如果检测到 it 导入到标准输入输出
​      子进程
​        前面的 command 调用
​          克隆隔离进程
​          进入进程,调用 自己
​          调用 init 初始化
​      初始化
​        容器内执行
​          挂载 /proc
​          替换 1 进程？

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