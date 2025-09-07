## 演示通过 etcd 加载业务配置

* 创建项目:  
trpc create -p ./pb/helloworld.proto -o . --mod=biz_config_etcd_demo -f

* 配置etcd 插件配置：
在 trpc_go.yaml 文件中新增配置项：

```
plugins:  # Plugin configuration.
  config:
    etcd:
      endpoints:
        - localhost:2379
      dialtimeout: 5s
```

* 定义 业务配置插件和配置和etcd 交互：
在文件 biz_config/biz_config_logic.go 文件中增加内容

* 在main.go 中增加对配置文件使用：

```
 bizconfig.WatchBizConfig()

 getBizConfig()
```


* 参考文档： https://github.com/trpc-group/trpc-go/blob/main/config/README.zh_CN.md 
* 参考文档： https://github.com/trpc-group/trpc-go/blob/main/docs/developer_guide/develop_plugins/config.zh_CN.md
