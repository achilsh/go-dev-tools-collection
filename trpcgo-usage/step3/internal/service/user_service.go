package service

import (
	"context"
	"step3/pb/inner/user"
	pb "step3/pb/inner/user"
)

// UserServiceImpl 实现UserService接口
type UserServiceImpl struct{}

// NewUserService 创建用户服务实例
func NewUserService() *UserServiceImpl {
	return &UserServiceImpl{}
}

func (u *UserServiceImpl) GetUserInfo(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	return &user.GetUserResponse{
		UserId:   req.UserId,
		Username: "test_user_" + req.UserId,
		Age:      25,
		Email:    "user_" + req.UserId + "@example.com",
	}, nil
}
