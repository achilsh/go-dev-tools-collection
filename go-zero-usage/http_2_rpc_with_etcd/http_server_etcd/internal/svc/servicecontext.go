package svc

import (
	"http_server_etcd/internal/config"
	pb "rpc_server_etcd/rpc_server_etcd"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config config.Config
	// 增加访问rpc的 client
	RpcClient pb.RpcServerEtcdClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 创建一个连接
	conn := zrpc.MustNewClient(c.ClientEtcd)

	return &ServiceContext{
		Config: c,
		// 初始化 rpc client 对象
		RpcClient: pb.NewRpcServerEtcdClient(
			conn.Conn(),
		),
	}
}
