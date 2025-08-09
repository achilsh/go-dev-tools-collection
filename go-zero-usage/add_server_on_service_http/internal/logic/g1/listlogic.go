package g1

import (
	"context"

	"add_server_on_service_http/internal/svc"
	"add_server_on_service_http/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListLogic {
	return &ListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListLogic) List() (resp []types.ListItem, err error) {
	// todo: add your logic here and delete this line

	return
}
