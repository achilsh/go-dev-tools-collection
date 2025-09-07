## 本项目演示 trpc-go 的 filter 使用

* trpc-go 的 filter 能实现的功能： 接口鉴权，日志记录，监控上报，服务端熔断限流器，敏感数据脱敏模块，服务端 panic 自动捕获插件，web referer 验证，参数自动校验插件。

* 如何实现 trpc-go 的filter: configure NewServer with ServerOption server.WithFilter()：
trpc.NewServer(server.WithFilter(demoFilter)); 客户端： client.WithFilter(demoFilter)
or configure trpc_go.yaml with sever filter.

* 已经的filter 有哪些？

* 通过 trpc_go.yaml配置来实现 自定义业务的filter:

1. 在 trpc-go.yaml 配置文件配置 filter 项。

1. 定义plugin，从配置文件加载 filter 的配置项
2. 使用 配置项的 插件名称 来 向filter 工厂注册服务端或者客户端 filter 函数

* 参考开源的filter: <https://github.com/trpc-ecosystem/go-filter/blob/main/README.zh_CN.md>
* 参考示例： <https://github.com/trpc-group/trpc-go/blob/main/examples/features/filter/README.md>
* filter原理介绍： <https://github.com/trpc-group/trpc-go/blob/main/filter/README.zh_CN.md>

* 生成一个项目：
trpc create -p ./pb/helloworld.proto -o . --mod=filter_demo -f

* 增加 filter plugin 的源码： biz_filter_logic 目录
* 在 服务端的 main.go 增加：
import (
    _"filter_demo/biz_filter_logic"
)

* 在服务端的trpc_go.yaml 配置文件 新增插件配置：

```
plugins:  # Plugin configuration.
  # 自定义filter 配置插件
  filterTypeOne:
    filterNameOne:
        abc: 123 
        xyz: achilsh 
```

在服务端的 server 上配置 filter:

```
  filter:  # List of interceptors for all service handler functions.
    - simpledebuglog
    - recovery  # Intercept panics from business processing goroutines created by the framework.
    ### 加入定义的插件； 插件名
    - filterNameOne
```
