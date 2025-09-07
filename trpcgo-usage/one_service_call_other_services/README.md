## 本项目演示 上游服务调用下游下游服务

* 先实现 下游服务，下游服务端提供多个 服务， 比如： second_service
生成 rpc 命令：
trpc create -p ./pb/second_server.proto -o . --mod=one_service_call_other_services -f
进入 stub目录后修改 go.mod 的 go 版本为 1.24； 运行 go mod tidy

* 在实现上游服务， 上游服务调用下游服务， 生成上午有服务 first_service
生成 rpc 命令：
trpc create -p ./pb/helloworld.proto -o . --mod=first_service_demo -f
进入 stub目录后修改 go.mod 的 go 版本为 1.24； 运行 go mod tidy

* 修改上游服务的函数内部处理， student_service.go 内部调用下游服务 call_second_service().
* 同时修改上游服务的trpc_go.yaml文件，新增下游客户的配置：

```
 ## 增加内部服务调用 下游客户端配置信息：
    - name: trpc.second.helloworld.Greeter  # Service name for the backend.
      namespace: Development  # Environment for the backend.
      network: tcp  # Network type for the backend: tcp or udp (configuration takes priority).
      protocol: trpc  # Application layer protocol: trpc or http.
      target: ip://127.0.0.1:8000  # Service address for requests.
      timeout: 1000   # Maximum processing time for requests.
    - name: trpc.second.helloworld.Hello  # Service name for the backend.
      namespace: Development  # Environment for the backend.
      network: tcp  # Network type for the backend: tcp or udp (configuration takes priority).
      protocol: trpc  # Application layer protocol: trpc or http.
      target: ip://127.0.0.1:8001  # Service address for requests.
      timeout: 1000   # Maximum processing time for requests.
```

* 分别启动下游服务；和 上游服务， 最后启动 上游服务的客户端程序，把整个链路走通：
上游服务客户端---> 上游服务 ----> 下游服务客户端（嵌入到上游服务内部） ----> 下游服务。



* 参考服务端开发： https://github.com/trpc-group/trpc-go/blob/main/docs/user_guide/server/overview.zh_CN.md
* 参考客户端开发： https://github.com/trpc-group/trpc-go/blob/main/docs/user_guide/client/overview.zh_CN.md
