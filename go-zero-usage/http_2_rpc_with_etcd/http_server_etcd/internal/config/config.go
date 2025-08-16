package config

import (
	"github.com/zeromicro/go-zero/core/discov"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	// 新增 http server 自身节点向etcd上报的etcd节点配置：
	RegisterEtcd discov.EtcdConf

	// 新增访问rpc服务的client配置；目前采用etcd方式
	ClientEtcd zrpc.RpcClientConf
}
