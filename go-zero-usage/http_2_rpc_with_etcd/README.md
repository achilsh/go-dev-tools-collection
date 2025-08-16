## 介绍 go-zero http call rpc by etcd

### etcd 本地安装和管理端界面安装

* etcd docker-compose 环境安装

```
进入 etcd_install/  运行 docker-compose up -d

```

* 安装 etcd管理端软件： Etcd Workbench

### 启动一个Http server,带有 etcd注册机制

* 使用go-zero第三方一个扩展：
<https://github.com/zeromicro/zero-contrib/blob/main/rest/registry/etcd/README.md>

修改配置文件，增加向etcd的注册, 修改 etc/httpserveretcd-api.yaml， 增加注册的etcd配置

```
RegisterEtcd:
  Hosts:
    - 127.0.0.1:2379
  Key: http_server.api.1

```

修改 配置文件代码，修改文件 config/config.go

```
type Config struct {
 rest.RestConf
 // 新增 http server 自身节点向etcd上报的etcd节点配置：
 RegisterEtcd discov.EtcdConf
}
```

启动 http 服务， 在etcd的管理端界面查看 上报的节点信息：
![alt text](http_register_etcd.png)
