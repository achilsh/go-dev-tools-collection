package g1

import (
	"context"

	"add_server_on_service_http/internal/svc"
	"add_server_on_service_http/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLogic {
	return &UpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateLogic) Update(req *types.UpdateReq) error {
	// todo: add your logic here and delete this line
	l.Infof("update call, req: %+v", *req)

	return nil
}
