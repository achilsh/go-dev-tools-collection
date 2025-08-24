package logic

import (
	"context"

	"http_gateway_demo/upstreams-server/userservice/internal/svc"
	"http_gateway_demo/upstreams-server/userservice/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserserviceUserIdLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserserviceUserIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserserviceUserIdLogic {
	return &UserserviceUserIdLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserserviceUserIdLogic) UserserviceUserId(req *types.Request) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line
	v := l.ctx.Value("x-abc")
	vs, ok := v.(string)
	if ok {
		l.Logger.Infof("get header data: %v", v)
	}
	resp = &types.Response{
		Message: "this is user id response: " + vs,
	}
	return
}
