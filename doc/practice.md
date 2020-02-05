[TOC]

### 7 高级实践

####   使用mydocker 创建一个可访问的nginx容器

​    获取nginx tar包
​    构建自己的nginx镜像
​    运行mynginx容器

####   使用mydocker 创建一个flask + redis的计数器

​    创建网桥
​    创建redis容器
​    制作flask镜像
​    创建myflask容器

####   runC

​    简介
​      OCI 组织
​      支持 namespace
​      支持 linux 提供的安全特性

#####     OCI 标准包（bundle）

​      config.json
​      rootfs
​    config.json
​      oci 版本号
​      root
​      下面应该都是配置
​    mounts
​      挂载点位置
​      类型
​      设备名
​      额外信息
​    process
​      容器进程信息
​      终端
​      控制台ui
​      cwd 可执行文件目录
​      env 环境变量
​      arg 参数
​      capabilities
​      rlimits
​    user
​      用户信息
​      uid
​      gid
​      additionalGids
​    hostname
​      主机名
​    platform
​      os 系统类型
​      arch 系统架构
​    钩子（Hook）
​      prestart
​      poststart
​      poststop

####   runC 创建容器流程

​    代码
​    流程
​      读取配置文件
​      设置 rootFileSystem
​      使用 factory 创建容器
​      创建容器初始化进程 process
​      设置容器输出管道 pipes
​      执行 Container.start() 启动容器
​      回调 init 方法重新 初始化容器
​      runc 父进程等待子进程初始化成功后退出

####   Docker containerd 项目介绍

​    每个 containerd 负责一台机器
​    架构
​      Distribution
​        拉镜像
​      Bundle
​        管理磁盘镜像
​      Runtime
​        创建容器，管理容器

#####     特性和路线图

​      支持 OCI镜像
​      支持 OCI 运行时 runC
​      支持 pull/push 镜像
​      容器运行时和生命周期管理
​      网络原语： 创建/修改/删除接口
​      容器加入已有 Network Nmesapces
​      全局镜像共享
​    containerd和Docker 之间的关系
​      Docker 功能多一点
​      Docker 可以构建镜像
​      containerd 提供 API 偏底层
​    containerd、OCI和runC之间的关系
​      OCI 规范
​      runC 实现
​      containerd 使用 runC，下载镜像，管理网络
​    containerd和容器编排系统的关系
​      被调用 和 docker 类似

####   Kubernetes CRI容器引擎

​    什么是CRI
​      容器运行时接口
​      包含
​        Protocol Buffers
​        gRPC API
​        运行库支持

#####       概览

​        kubelet 通过 gRPC 与 CRI shim 通信
​        CRI shim 通过 Unix Socket 启动 gPRC server 提供容器运行时
​        kubelet 做为 gPRC client 通过 Unix Socket 与 CRI shim 通信
​        gPRC server 使用 protocol buffrs 提供两类 gPRC server
​        ImageService拉取，删除，查询镜像的RPC调用功能
​        RuntimeService容器创建，修改，销毁及交互操作

#####       接口实现

​        PodSandbox
​        container
​        不同运行时的 PodSandbox 实现不一样
​        Hypervisor 一组虚拟机， docker 命名空间

#####       RuntimeService

​        创建 pod 前调用
​        其为 pod 创建一个 PodSandbox
​        初始化pod网络，分配IP，激活沙箱
​        然后 kubelet 再进行各种操作 建启停删
​        删除时，也会调用其
​        还是容器进行交互的接口

#####       ImageService

​        镜像拉取，查看，移除
​        不包括构建
​      LogService
​        stdout/stderr 日志处理规范？
​    为什么需要CRI
​      CRI 之前要通过 high-level 接口集成，要懂k8s
​      使用可插拔容器运行时
​    为什么CRI是接口且是基于容器的而不是基于Pod的
​      负担重
​      pod还在发展
​    如何使用CRI
​      添加标记
​    CRI的目标
​    已知的问题