package logic

import (
	"context"

	"rpc_demo_server2/internal/svc"
	"rpc_demo_server2/rpc_demo_server2"

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

func (l *PingLogic) Ping(in *rpc_demo_server2.Request) (*rpc_demo_server2.Response, error) {
	// todo: add your logic here and delete this line

	l.Logger.Infof("req_v3: %v", in.Ping)
	return &rpc_demo_server2.Response{
		Pong: "pong response 3",
	}, nil
}
