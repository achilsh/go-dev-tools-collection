package g1

import (
	"context"

	"add_server_on_service_http/internal/svc"
	"add_server_on_service_http/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FormExampleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFormExampleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FormExampleLogic {
	return &FormExampleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FormExampleLogic) FormExample(req *types.FormExampleReq) error {
	// todo: add your logic here and delete this line

	return nil
}
