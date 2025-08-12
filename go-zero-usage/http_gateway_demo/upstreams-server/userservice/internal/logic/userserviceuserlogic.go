package logic

import (
	"context"

	"http_gateway_demo/upstreams-server/userservice/internal/svc"
	"http_gateway_demo/upstreams-server/userservice/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserserviceUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserserviceUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserserviceUserLogic {
	return &UserserviceUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserserviceUserLogic) UserserviceUser() (resp *types.Response, err error) {
	// todo: add your logic here and delete this line

	return
}
