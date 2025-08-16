package logic

import (
	"context"

	"hump_style/internal/svc"
	"hump_style/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type Hump_styleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewHump_styleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *Hump_styleLogic {
	return &Hump_styleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *Hump_styleLogic) Hump_style(req *types.Request) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line

	return
}
