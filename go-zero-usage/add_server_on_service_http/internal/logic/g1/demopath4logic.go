package g1

import (
	"context"

	"add_server_on_service_http/internal/svc"
	"add_server_on_service_http/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DemoPath4Logic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDemoPath4Logic(ctx context.Context, svcCtx *svc.ServiceContext) *DemoPath4Logic {
	return &DemoPath4Logic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DemoPath4Logic) DemoPath4(req *types.DemoPath4Req) (resp *types.DemoResp, err error) {
	// todo: add your logic here and delete this line

	return
}
