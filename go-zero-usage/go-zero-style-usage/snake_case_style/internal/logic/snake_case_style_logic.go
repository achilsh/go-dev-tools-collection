package logic

import (
	"context"

	"snake_case_style/internal/svc"
	"snake_case_style/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type Snake_case_styleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSnake_case_styleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *Snake_case_styleLogic {
	return &Snake_case_styleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *Snake_case_styleLogic) Snake_case_style(req *types.Request) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line

	return
}
