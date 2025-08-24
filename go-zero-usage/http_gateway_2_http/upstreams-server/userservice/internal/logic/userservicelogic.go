package logic

import (
	"context"
	"fmt"

	"http_gateway_demo/upstreams-server/userservice/internal/svc"
	"http_gateway_demo/upstreams-server/userservice/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserserviceLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserserviceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserserviceLogic {
	return &UserserviceLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserserviceLogic) Userservice(req *types.Request) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line
	l.Debugf("request: %+v", *req)
	resp = &types.Response{
		Message: fmt.Sprintf("%v", req.ID),
	}
	return resp, nil
}
