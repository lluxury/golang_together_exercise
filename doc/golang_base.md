#### 编译

非Linux下编写 Linux api相关，要加标记

```
// +build linux
```



编译要指定参数，然后把产物传到 linux 环境下执行，也并不方便

```
GOOS=linux go build ns-proc.go
```



符号链接内容

```
readlink /proc/self/ns/uts  # 测试
```

