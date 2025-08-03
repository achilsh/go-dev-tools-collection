package service

import (
	"context"
	"fmt"

	pb "step3/pb/api/helloworld"

	"step3/pb/inner/user"

	"trpc.group/trpc-go/trpc-go/client"
	"trpc.group/trpc-go/trpc-go/log"
)

// HelloworldService is used to implement pb.GreeterService.
type HelloworldService struct {
	// UnimplementedHelloworld is the unimplemented greeter service server
	pb.UnimplementedHelloworld

	// 集成调用下游的 trpc 客户端
	userClient user.UserServiceClientProxy
}

func NewHelloworldService() *HelloworldService {
	return &HelloworldService{
		userClient: user.NewUserServiceClientProxy(
			client.WithTarget("ip://127.0.0.1:8001"),
			client.WithProtocol("trpc"),
		),
	}
}

// SayHello RESTful HTTP接口：通过用户ID获取问候信息
//
//	会内部调用UserService获取用户信息
//
// 映射 url: /v1/greeter/hello?user_id=1231
func (g HelloworldService) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	log.InfoContextf(ctx, "[restful] Received SayHello request with req: %v", req)
	inresp, err := g.userClient.GetUserInfo(ctx, &user.GetUserRequest{
		UserId: "inner_req: " + req.GetUserId(),
	})

	if err != nil {
		// handle request
		rsp := &pb.HelloResponse{
			Message: "[restful] SayHello Hello: " + fmt.Sprintf("inner err: %v", err.Error()),
		}
		return rsp, nil
	}

	rsp := &pb.HelloResponse{
		Message: "[restful] SayHello Hello: " + ", inner response: " + inresp.Email,
	}
	return rsp, nil
}

// CreateUser RESTful HTTP接口：创建用户（内部调用UserService的扩展接口）
func (g HelloworldService) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	log.InfoContextf(ctx, "[restful] Received CreateUser request with req: %v", req)
	// handle request

	rsp := &pb.CreateUserResponse{
		Message: "[restful] SayHello Hello " + req.Username,
	}
	return rsp, nil
}
