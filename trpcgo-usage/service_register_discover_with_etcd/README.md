## 本项目演示 trpc-go 使用 etcd 做服务注册和发现，演示上下游通过服务名访问

* 1. 演示 user 服务（user_demo） 调用 student, 其中 student 为多个服务实例： student_demo1/student_demo2

### 被调用方 student_demo 源码创建

* 生成源码命令：
trpc create -p ./pb/helloworld.proto  --mod=student_demo1 -f --alias -o .

* 增加 etcd 的服务注册插件：
在 main.go 中 import () etcd 服务注册的插件

```
import (
     //  引入 etcd 的服务注册发现插件
 _ "trpc.group/trpc-go/trpc-naming-etcd"
 _ "trpc.group/trpc-go/trpc-naming-etcd/registry"
)
```

* 修改 trpc_go.yaml 文件：
服务端配置

```

  - name: trpc.student.demo.HelloWorldService  # Route name for the service.
    # ip: 127.0.0.1  # Service listening IP address, can use placeholder ${ip}. Use either ip or nic, ip takes priority.
    nic: eth0
    port: 8000  # Service listening port, can use placeholder ${port}.
    network: tcp  # Network listening type: tcp or udp.
    protocol: trpc  # Application layer protocol: trpc or http.
    timeout: 1000  # Maximum processing time for requests in milliseconds.
  - name: http.student.demo.HelloWorldService  # Service name for the backend.
    nic: eth0
    port: 7000  # Service listening port, can use placeholder ${port}.
    network: tcp  # Network listening type: tcp or udp.
    protocol: http  # Application layer protocol: trpc or http.
    timeout: 1000  # Maximum processing time for requests in milliseconds.
```

新增插件：

```
plugins:  # Plugin configuration.
  registry:
    etcd:
      address: 127.0.0.1:2379
      timeout: 5
      service:
      - name: trpc.student.demo.HelloWorldService  ### 向etcd 注册的服务名，和上面 server.service.name 的值保持一致
        ttl: 10
        metadata:
          tags: helloworld
      - name: http.student.demo.HelloWorldService  ### 向etcd 注册的服务名，和上面 server.service.name 的值保持一致
        ttl: 10
        metadata:
          tags: helloworld
```

### 被调用方 student_demo2 源码创建

trpc create -p ./pb/helloworld.proto  --mod=student_demo2 -f --alias -o .

* 增加 etcd 的服务注册插件：
在 main.go 中 import () etcd 服务注册的插件

```
import (
     //  引入 etcd 的服务注册发现插件
 _ "trpc.group/trpc-go/trpc-naming-etcd"
 _ "trpc.group/trpc-go/trpc-naming-etcd/registry"
)
```

* 修改 trpc_go.yaml 文件：
服务端修改配置

```
  service:  # Services provided by the business, can have multiple.
  - name: trpc.student.demo.HelloWorldService  # Route name for the service.
    # ip: 127.0.0.1  # Service listening IP address, can use placeholder ${ip}. Use either ip or nic, ip takes priority.
    nic: eth0
    port: 8001  # Service listening port, can use placeholder ${port}.
    network: tcp  # Network listening type: tcp or udp.
    protocol: trpc  # Application layer protocol: trpc or http.
    timeout: 1000  # Maximum processing time for requests in milliseconds.
  - name: http.student.demo.HelloWorldService  # Route name for the service.
    # ip: 127.0.0.1  # Service listening IP address, can use placeholder ${ip}. Use either ip or nic, ip takes priority.
    nic: eth0
    port: 7001  # Service listening port, can use placeholder ${port}.
    network: tcp  # Network listening type: tcp or udp.
    protocol: http  # Application layer protocol: trpc or http.
    timeout: 1000  # Maximum processing time for requests in milliseconds.

```

新增插件:

```
plugins:  # Plugin configuration.
  registry:
    etcd:
      address: 127.0.0.1:2379
      timeout: 5
      service:
      - name: trpc.student.demo.HelloWorldService  ### 向etcd 注册的服务名，和上面 server.service.name 的值保持一致
        ttl: 10
        metadata:
          tags: helloworld
      - name: http.student.demo.HelloWorldService  ### 向etcd 注册的服务名，和上面 server.service.name 的值保持一致
        ttl: 10
        metadata:
          tags: helloworld

```

### 调用方 实例创建

* 创建项目命令：
trpc create -p ./pb/helloworld.proto  --mod=user_demo -f --alias -o .

* 增加 etcd 的服务注册插件： 在 main.go 中 import () etcd 服务注册的插件

```
import (
 //  集成 etcd 服务发现的插件
 _ "trpc.group/trpc-go/trpc-naming-etcd"
 _ "trpc.group/trpc-go/trpc-naming-etcd/registry"
)
```

* 在调用方集成 服务发现的配置， 修改 trpc_go.yaml 配置：

```
client: #客户端调用的后端配置
  service: #针对单个后端的配置
    - callee: trpc.test.helloworld.Greeter         #后端服务协议文件的service name, 如何callee和下面的name一样，那只需要配置一个即可
      target: etcd://trpc.test.helloworld.Greeter              #后端服务地址 etcd
      network: tcp                                 #后端服务的网络类型 tcp udp
      protocol: http                              #应用层协议 trpc http
      timeout: 10000                               #请求最长处理时间
```
