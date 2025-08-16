package logic

import (
	"context"

	"http_server_etcd/internal/svc"
	"http_server_etcd/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type Http_server_etcdLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewHttp_server_etcdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *Http_server_etcdLogic {
	return &Http_server_etcdLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *Http_server_etcdLogic) Http_server_etcd(req *types.Request) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line

	return
}
