package logic

import (
	"context"

	"rpc_server_etcd/internal/svc"
	"rpc_server_etcd/rpc_server_etcd"

	"github.com/zeromicro/go-zero/core/logx"
)

type PingLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PingLogic {
	return &PingLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PingLogic) Ping(in *rpc_server_etcd.Request) (*rpc_server_etcd.Response, error) {
	// todo: add your logic here and delete this line
	return &rpc_server_etcd.Response{
		Pong: "this is message pong.----------------------------",
	}, nil
}
