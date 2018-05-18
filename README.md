### Pierce

#### 依赖
```
# go version 1.9
github.com/golang/glog
github.com/golang/protobuf
```

### 功能
内网穿透。

举个例子，有一台内网的机器A，没有公网ip，但是可以访问公网。有一台云服务器B，有公网。

那么可以实现端口转发，将A机器的端口映射到B机器上。

### 原理
机器A主动连接机器B，建立MainPipe，断开时自动重连，并用心跳包保持连接。

当机器B收到tcp请求时，通过MainPipe告知机器A，然后机器A与对应的InnerIp:InnerPort发tcp请求，并同时给机器Ｂ的OuterIp:OuterPort发tcp请求，建立Pipe，之后这次请求的所有数据全部走这条Pipe。

### 安装
```
# 设置好$GOPATH, $GOROOT, $GOBIN
# 用go get github.com/qwb2333/Pierce 下载代码
./install.py
```

### 运行
cd $GOBIN
./inner ../config/inner.properties

./outer ../config/outer.properties ../config/outer_dot.conf

记得注意配置文件的路径

### 注意事项
inner.properties中的OuterIp填机器B的公网ip, outer.properties的OuterIp通常填0.0.0.0,否则收不到本机外的请求。

outer_dot.conf中如果OuterIp如果填127.0.0.1则只能在本机访问，如果要其他地方能访问，需要填0.0.0.0。
