package main

import (
	"context"
	"fmt"

	pb "github.com/achilsh/go-dev-tools-collection/trpcgo-usage/step2/helloworld"
	trpc "trpc.group/trpc-go/trpc-go"
	"trpc.group/trpc-go/trpc-go/log"
)

func main() {
	// init trpc server
	server := trpc.NewServer()
	// Register the greeter service with the server
	pb.RegisterGreeterService(server, new(greeterService))
	// Run the server
	if err := server.Serve(); err != nil {
		log.Fatal(err)
	}
}

// greeterService is used to implement pb.GreeterService.
type greeterService struct {
	// unimplementedGreeterServiceServer is the unimplemented greeter service server
	pb.UnimplementedGreeter
}

func (g greeterService) SayHello(
	ctx context.Context,
	req *pb.HelloRequest,
) (rsp *pb.HelloReply, err error) {

	log.InfoContextf(ctx, "[restful] Received SayHello request with req: %v", req)
	// handle request
	rsp = &pb.HelloReply{
		Message: "[restful] SayHello Hello " + req.Name,
	}
	return rsp, nil
}

func (g greeterService) DoPostMessage(
	ctx context.Context,
	req *pb.PostMessageRequest,
) (rsp *pb.PostMessageResponse, err error) {
	log.InfoContextf(ctx, "[restful] Received Message request with req: %v", req)
	// handle request
	rsp = &pb.PostMessageResponse{
		Data: fmt.Sprintf("[restful] Message name:%s, age:%d",
			req.GetName(), req.GetAge()),
		Code: 99,
	}
	return rsp, nil

}

// DoGetMessage GET 请求：/v1/get/info?age=1&address=shenzhen
//
//	查询参数 age 和 address 自动映射到 GetMessageRequest 的同名字段
func (g greeterService) DoGetMessage(
	ctx context.Context,
	req *pb.GetMessageRequest,
) (rsp *pb.GetMessageResponse, err error) {
	log.InfoContextf(ctx, "[restful] Received Message request with req: %v", req)
	// handle request
	rsp = &pb.GetMessageResponse{
		Message: fmt.Sprintf("[restful] Message address:%s, age:%d",
			req.GetAddress(), req.GetAge()),
		Code: 100,
	}
	return rsp, nil
}
