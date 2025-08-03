## 通过trpc-go 实现一个对外的 restful http server, 对内实现一个 trpc server; http server 调用 内部 trpc server

### 1. 定义 pb 文件，包括对外和对内服务调用的文件

* 根据 user.proto  生成源文件命令：

```
trpc create -d pb/inner -p user.proto --mock=false  --nogomod --rpconly --gotag --alias -f -o pb/inner/user/
```

* 根据 helloworld.proto 生成源码：

```
trpc create  -d pb/  -p api/helloworld.proto  --mock=false  --nogomod --rpconly --gotag --alias -f -o pb/api/helloworld/
```

* 增加 内部服务的主函数入口： cmd/user_server/main.go
* 增加 内部服务对外提供 trpc 的接口： internal/service/user_service.go

* 增加 http 服务处理，实现 http 接收 restful 请求的处理：internal/service/helloworld.go, 包括集成内部服务的 client，在收到外部请求时，通过client调用内部服务。

* 增加restful http server 的配置文件： config/helloworld.yaml
* 增加 内部服务的配置文件： config/user.yaml

* 启动内部服务 go run cmd/user_server/main.go -conf config/user.yaml
* 启动http api服务 go run cmd/http_server/main.go -conf ./config/helloworld.yaml

* 在浏览器输入： <http://127.0.0.1:9093/v1/greeter/hello?user_id=1>
* 查看结果

## TOOD

* http server 内部调用 下游服务的client 配置。
* 两个服务的其他配置
