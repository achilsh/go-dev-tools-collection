package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	//  http 服务端配置
	rest.RestConf

	//  rpc client 的配置定义
	DirectConnetClientCfg zrpc.RpcClientConf
}
