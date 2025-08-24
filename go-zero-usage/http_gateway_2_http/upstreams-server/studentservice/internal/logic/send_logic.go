package logic

import (
	"context"

	"http_gateway_demo/upstreams-server/studentservice/internal/svc"
	"http_gateway_demo/upstreams-server/studentservice/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendLogic {
	return &SendLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SendLogic) Send(req *types.SendReq) (resp *types.SendResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
