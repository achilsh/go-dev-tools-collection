package logic

import (
	"context"

	"rpc_demo_server1/internal/svc"
	"rpc_demo_server1/rpc_demo_server1"

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

func (l *PingLogic) Ping(in *rpc_demo_server1.Request) (*rpc_demo_server1.Response, error) {
	// todo: add your logic here and delete this line

	return &rpc_demo_server1.Response{
		Pong: "111111111111111",
	}, nil
}
