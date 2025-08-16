package logic

import (
	"context"

	"low_case_style/internal/svc"
	"low_case_style/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type Low_case_styleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLow_case_styleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *Low_case_styleLogic {
	return &Low_case_styleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *Low_case_styleLogic) Low_case_style(req *types.Request) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line

	return
}
