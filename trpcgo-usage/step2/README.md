# 使用 trpc-go 实现一个 http 服务 和 trpc 服务

* 使用 trpc-go 对外提供http 服务方式目前存在两种，1： 集成http 标准库，2： 使用 trpc-go restful 插件。两种均可享受 TRPC 框架的服务治理能力（如监控、日志、中间件），无需担心核心功能缺失。

* 下面使用 trpc restful 特性来实现 http ---> trpc 内部服务

* restful 使用 示例参考： <https://github.com/trpc-group/trpc-go/tree/v1.0.3/examples/features> 目录中的 restful。

* restful 协议定义文档： <https://github.com/trpc-group/trpc-go/tree/main/restful>

## 根据 Pb helloworld.proto 创建单一的http server

* 创建源码：

```
trpc create -d pb  -p helloworld.proto --mock=false  --nogomod --rpconly --gotag --alias -f -o helloworld/
```
