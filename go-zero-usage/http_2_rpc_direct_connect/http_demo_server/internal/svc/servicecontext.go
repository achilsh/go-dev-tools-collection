package svc

import (
	"http_demo_server/internal/config"

	pb "rpc_demo_server1/rpc_demo_server1"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config config.Config
	// 增加客户端句柄信息
	DirectConnClient pb.RpcDemoServer1Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	retCtx := &ServiceContext{
		Config: c,
	}

	//创建 rpc client并初始化
	retCtx.DirectConnClient = pb.NewRpcDemoServer1Client(zrpc.MustNewClient(c.DirectConnetClientCfg).Conn())
	return retCtx
}
