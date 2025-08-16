package svc

import (
	"http_demo_server/internal/config"

	pb "rpc_demo_server1/rpc_demo_server1"
	mp2 "rpc_demo_server2/rpc_demo_server2"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config config.Config
	// 增加客户端句柄信息
	DirectConnClient pb.RpcDemoServer1Client
	// 多服务端的节点直连
	MoreCliConnClients mp2.RpcDemoServer2Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	retCtx := &ServiceContext{
		Config: c,
	}

	//创建 单一rpc client并初始化
	retCtx.DirectConnClient = pb.NewRpcDemoServer1Client(zrpc.MustNewClient(c.DirectConnetClientCfg).Conn())
	// 初始化 多个 rpc client
	retCtx.MoreCliConnClients = mp2.NewRpcDemoServer2Client(zrpc.MustNewClient(c.MoreNodeConnClientCfg).Conn())

	return retCtx
}
