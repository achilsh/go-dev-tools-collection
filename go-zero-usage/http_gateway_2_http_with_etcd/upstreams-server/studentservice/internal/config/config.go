package config

import (
	"github.com/zeromicro/go-zero/core/discov"
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf

	// 新增 http server 自身节点向etcd上报的etcd节点配置：
	Etcd discov.EtcdConf
}
