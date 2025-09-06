## 如何使用trpc 框架配置

* 参考文档： <https://github.com/trpc-group/trpc-go/blob/main/docs/user_guide/framework_conf.zh_CN.md>

* 使用配置文件原则：

```
使用框架配置文件，trpc.NewServer() 在启动时，会先解析框架配置文件，自动初始化所有配置好的插件，并启动服务。建议其他初始化逻辑都放在 trpc.NewServer() 之后，以确保框架功能已经初始化完成。tRPC-Go 的默认框架配置文件名称是trpc_go.yaml，默认路径为当前程序启动的工作路径，可通过 -conf 命令行参数指定配置路径
```

trpc_go.yaml文件中有哪些内容，没想内容的含义可参考 上面文档介绍

* 下面是 trpc_go.yaml 内容的每一项示例， 参考：
<https://github.com/trpc-group/trpc-go/blob/main/testdata/trpc_go.yaml>
