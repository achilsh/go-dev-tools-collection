## trpc 插件使用介绍

* 创建一个 项目，生成源码：
trpc create -p ./pb/helloworld.proto -o . --mod=plugin_demo  -f

* 详细文档参考：1. <https://github.com/trpc-group/trpc-go/blob/main/plugin/README.zh_CN.md>

2. <https://github.com/trpc-group/trpc-go/blob/main/docs/basics_tutorial.zh_CN.md>

* 需要依赖配置进行加载的插件 的开发：
参考： <https://github.com/trpc-group/trpc-go/blob/main/examples/features/plugin/custom_plugin.go>
