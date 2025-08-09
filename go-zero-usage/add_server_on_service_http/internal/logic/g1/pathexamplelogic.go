package g1

import (
	"context"

	"add_server_on_service_http/internal/svc"
	"add_server_on_service_http/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PathExampleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPathExampleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PathExampleLogic {
	return &PathExampleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PathExampleLogic) PathExample(req *types.PathExampleReq) (resp *types.PathExampleResp, err error) {
	// todo: add your logic here and delete this line

	return
}
