package logic

import (
	"context"

	"demo/internal/svc"
	"demo/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type Demo_apiLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDemo_apiLogic(ctx context.Context, svcCtx *svc.ServiceContext) *Demo_apiLogic {
	return &Demo_apiLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *Demo_apiLogic) Demo_api(req *types.Request) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line

	return
}
