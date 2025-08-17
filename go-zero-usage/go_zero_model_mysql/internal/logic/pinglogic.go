package logic

import (
	"context"

	"go_zero_model_mysql/go_zero_model_mysql"
	"go_zero_model_mysql/internal/svc"

	"go_zero_model_mysql/model/mysql/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type PingLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PingLogic {
	return &PingLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PingLogic) Ping(in *go_zero_model_mysql.Request) (*go_zero_model_mysql.Response, error) {
	// todo: add your logic here and delete this line
	ret, err := l.svcCtx.UserModel.Insert(context.Background(), &user.User{
		Name: "this demo insert user",
		Type: 100,
	})
	if err != nil {
		l.Logger.Errorf("err insert: %v", err)
		return nil, err
	}
	l.Logger.Infof("insert data: %v", ret)
	return &go_zero_model_mysql.Response{}, nil
}

// 直接调用接口操作 db
func AddInsert(svcCtx *svc.ServiceContext, l logx.Logger) error {
	ret, err := svcCtx.UserModel.Insert(context.Background(), &user.User{
		Name: "this demo insert user",
		Type: 100,
	})
	if err != nil {
		l.Errorf("err insert: %v", err)
		return err
	}
	l.Infof("insert data: %v", ret)

	return nil

}
