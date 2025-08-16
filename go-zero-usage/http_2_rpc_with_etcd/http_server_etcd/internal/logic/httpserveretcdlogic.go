package logic

import (
	"context"
	"fmt"

	"http_server_etcd/internal/svc"
	"http_server_etcd/internal/types"

	pb "rpc_server_etcd/rpc_server_etcd"

	"github.com/zeromicro/go-zero/core/logx"
)

type Http_server_etcdLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewHttp_server_etcdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *Http_server_etcdLogic {
	return &Http_server_etcdLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *Http_server_etcdLogic) Http_server_etcd(req *types.Request) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line

	//通过 rpc client 调用下游服务
	rpcReq := &pb.Request{
		Ping: fmt.Sprintf("from http req: %v", req.Name),
	}

	ret, err := l.svcCtx.RpcClient.Ping(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}
	resp = &types.Response{
		Message: fmt.Sprintf("from rpc response: %v", ret.Pong),
	}
	return
}
