package logic

import (
	"context"

	"http_gateway_demo/upstreams-server/userservice/internal/svc"
	"http_gateway_demo/upstreams-server/userservice/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserServicePostLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserServicePostLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserServicePostLogic {
	return &UserServicePostLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserServicePostLogic) UserServicePost(req *types.PRequest) (resp *types.PResponse, err error) {
	// todo: add your logic here and delete this line

	v := l.ctx.Value("x-abc")
	l.Infof("x-abc: ", v)
	l.Infof(">>>>>>>>> req: %+v", *req)
	return
}
