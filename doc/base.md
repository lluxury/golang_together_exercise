[TOC]

### 2 基础技术



####   Linux Namespace 介绍

​    概念
​      UID 级别隔离

#####   系统调用

​    clone()创建**新进程**
​    unshare() 进程移出 Namespace
​    setns() 进程加入 Namespace

unshare 命令：

```bash
unshare -h
unshare -u /bin/sh  # 进入了新的名字空间 --uts
 # 类似 syscall.CLONE_NEWUTS
```



###### 名字空间梳理：

- Mount - 用于独立文件系统和挂载点 
- UTS -用于独立hostname和domainname 
- IPC - 用于独立进程之间的通信资源nterprocess communication (IPC) 
- PID - 用于独立PID号空间
- Network -用于间隔网络接口interfaces 
- User - 用于间隔UID/GID空间 
- **Cgroup** - 用于间隔cgroup root目录



exec.md

#####     UTS Namespace

```bash
# syscall.CLONE_NEWUTS

pstree -pl
  ...main(2603)─┬─sh(2606)
echo $$
readlink /proc/2603/ns/uts
readlink /proc/2606/ns/uts  # 父子进程不同

hostname -b bird 

# 新开窗口验证主机名
```



#####     IPC Namespace

​      ipcs -q
​      ipcmk -Q
​      运行代码
​      ipcs -q 确认新进程中无队列

消息队列



#####     PID Namespace

​      pstree -pl
​      echo $$
​      ps top 使用/proc 内容



#####     Mount Namespace

​      隔离进程看到的挂载点视图 类chroot
​       运行代码, ls /proc
​      **mount -t proc proc /proc**
​      ls /proc, ps-ef, 独立的进程内容



#####     User Namespace

​      id, 运行代码
​      id确认, 变成非root用户



#####     Network Namespace

​      ifconfig,运行代码
​      ifconfig确认 新空间无网络设备



####   Linux Cgroups介绍

#####     Linux Cgroups

######       cgroup

​        一组进程

######       subsystem

​        安装 

​		**apt-get install cgroup-bin**
​        **lssubsys**

######       hierarchy

​        把cgroup串树状,属性继承
​       

关系：

- 系统创建hierarchy后,所有进程加入根节点

- 一个subsystem 只能附在一个hierarchy上  

- 一个进程可以存在多个cgroup中,只要cgroup属于不同hierarchy 

- 推出 hierarchy 中,进程是唯一

- fork 出的子进程默认和父进程同 cgroup 可以移出

  

#####     hierarchy实践

######      kernel接口 含义

​      cgroup.clone_children
​      cgroup.procs
​      cgroup.sane_behavior
​      notify_on_release
​      release_agent
​      tasks



######       建立cgroup树

​        mkdir cgroup-test
​        **sudo mount -t cgroup -o none,xxx**
​        ls ./cgroup-test/



######       扩展子节点

​        apt-get install tree
​        sudo mkdir cgroup-1
​        sudo mkdir cgroup-2
​        tree



######       添加移动进程

​        cd cgroup-1/
​        echo $$
​        sudo sh -c "echo $$" >> tasks
​        **cat /proc/4252/cgroup**
​        确认当前进程,记录在cgroup-1中



######       限制cgroup中进程资源

​        cd /sys/fs/cgroup/memory
​        stress --vm-bytes 200m --vm-keep -m 1
​        压测,top确认内存占用, RES
​        test-limit-memory 
​        memory.limit_in_bytes
​        sudo sh -c "echo $$ > tasks"
​        建立子系统,写入限制大小,tasks写入进程号

通过 挂载 cgrop 建立, 写入 tasks 转移进程



#####     Docker是如何使用Cgroups的

​      apt install docker.io
​      docker run -itd -m 128m ubuntu
​      cd /sys/fs/cgroup/memory/docker/ID
​      cat memory.limit_in_bytes
​      cat memory.usage_in_bytes



#####     用Go语言实现通过cgroup限制容器的资源

​      判断执行文件路径
​        如果符合 /proc/self/exe 就继续
​          这个后面讲
​        打印 进程 pid
​        发起一个 stress 进程
​          和上节使用的类似
​        Subtopic
​        定义控制台的输入输出及错误
​      定义控制台的输入输出及错误
​      获取 cmd的执行状态
​        如果 有err 则报错
​        否则打印fork出来进程映射在外部命名空间的pid
​          cmd.Process.Pid
​        在默认Hierarchy上创建cgroup
​          os.Mkdir(path.Join())
​        把容器进程加入新建 cgroup
​          ioutil.WriteFile(path.Join())
​        在对应位置写入限制数据
​          ioutil.WriteFile(path.Join())



####   Union File System

#####     什么是Union File System

​      多个文件系统叠加
​      写时复制 CoW
​      Knoppix 光盘linux



#####     AUFS

​      新的UnionFS

不同版本操作系统使用了不同的文件系统，遇到再记录



​    Docker是如何使用AUFS的

#####     image layer和AUFS

​      ls **/var/lib/docker/aufs/**
​      diff  放只读层image layer内容
​      layers  放堆栈layer的metadata 元数据
​      上面都是镜像相关的操作



#####     实践

​      编辑 Dockerfile
​      docker build -t changed-ubuntu .
​      docker images
​      docker history changed-ubuntu
​      查看执行的操作



#####     container layer和AUFS

​      read-only的init laye 和 read-write的layer
​      先确认下面2个位置没有文件,无容器运行
​      ll **/var/lib/docker/containers/**ID/  # 容器的 metadata 和 config 
​      ls /sys/fs/aufs/
​      cat **/sys/fs/aufs/si_**ead33945636649b3/*  # 容器权限
​      ls -a **/var/lib/docker/aufs/diff/**ID/tmp/   #2个初始层
​      删除文件  .wh.newfile
​      以上都是容器相关的操作位置



#####     自己动手写AUFS

​      构造aufs系统,修改文件,确认修改保存位置
​      mkdir mnt container-layer
​      准备4个目录,内有4个文件

主要用来通过 si* 看只读层的

```bash
sudo mount -t aufs -o dirs=./container-layer:./image-layer4:./image-layer3:./image-layer2:./image-layer1 none ./mnt

ls  /sys/fs/aufs/si_c2682246e4a009a4/*
```


​      修改 image-layer4
​      ls diff/container-layer/ 查找l4复制过来的文件

