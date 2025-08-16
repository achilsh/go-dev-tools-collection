package messagelogic

import (
	"context"

	"rpc_demo_1/internal/svc"
	"rpc_demo_1/pb/gen"

	"github.com/zeromicro/go-zero/core/logx"
)

type PongLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPongLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PongLogic {
	return &PongLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PongLogic) Pong(in *gen.SendMessageReq) (*gen.SendMessageResp, error) {
	// todo: add your logic here and delete this line

	return &gen.SendMessageResp{}, nil
}
