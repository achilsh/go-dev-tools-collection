## 演示加载 业务配置； 实现简单的本地业务配置文件的加载

* 创建一个业务trpc server：
trpc create -p ./pb/helloworld.proto -o . --mod=biz_config_demo -f

* 如何使用业务配置, 参考： <https://github.com/trpc-group/trpc-go/blob/main/config/README.zh_CN.md>

* 具体使用：

1. 在 main.go 文件中增加： loadBizConfig() 函数。
2. 然后 使用配置项中的内容
