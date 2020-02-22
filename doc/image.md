[TOC]

### 4 构造镜像

####   使用busybox创建容器

​    原因
​      可以看到父进程挂载点
​    busybox

#####       获取rootfs

​      docker export打包
​      tar 解压

```bash
docker pull busybox
docker run -d busybox top -b
docker export -o busybox.tar 31971ffadca9
mkdir busybox
tar xvf busybox.tar -C busybox/

```

运行镜像，导出，解压为footfs    



##### pivot_root

​      原理
​        改变当前的root文件系统
​        从新挂载
​        tmpfs 内存文件系统
​     

 代码

###### init.go

  setUpMount()
  pivotRoot()

  bind mount 成不同系统
  建目录 .pivot_root
  PivotRoot 移走root
  修改当前目录到根目录
  卸载老root
  删除临时挂载点



###### container_process.go

```go
  cmd.Dir = "/root/busybox"
```



####   使用AUFS包装busybox

​    原因
​      虽然换目录挂载了，还是会影响到宿主机里的现实目录
​      需要容器和镜像隔离
​      上一章的 AUFS基础内容

​    代码

###### container_process.go

  NewWorkSpace()
  CreateReadOnlyLayer()
  CreateWriteLayer()
  CreateMountPoint()
  更新NewParentProcess()

  DeleteWorkSpace()
  DeleteMountPoint()
  DeleteWriteLayer()

  PathExists()

挂载时常量变成了函数，卸载时 在 run.go 里面加一次调用

​    流程
​      创建工作目录
​      解压镜像包
​      确认解压成功
​      创建容器可写层
​      确认创建成功
​      创建 mnt 文件夹，其下挂载可写层和镜像
​      确认完成

####   实现volume数据卷

​    原因
​      容器退出，可写层内容删除？
​      原理
​        启动容器
​          创建只读层
​          创建读写层
​          创建挂载点，挂载以上两层
​        容器退出？
​          卸载 mnt
​          删除挂载点
​          删除读写层？
​        添加 -v 标签
​          挂盘后 判断 volume 值
​          不为空，解析字符串
​          验证通过则挂载
​          否则给出错误提示
​        挂载
​          建宿主机目录 /root/${parentUrl}
​          容器挂载点 /root/mnt/${containerUrl}
​          挂载
​    代码
​    流程
​      解析 Volume 参数，创建工作目录
​      抽取 宿主机 与 容器
​      创建文件，把宿主机目录挂载到容器内
​      检查创建成功
​    验证

####   实现简单镜像打包

​    原因
​      实现 docker commit
​      指定目录打包成镜像
​    代码
  小结