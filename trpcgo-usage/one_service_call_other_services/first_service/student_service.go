package main

import (
	"context"

	pb "hello_world_demo/my_hello_world"
)

type studentServiceImpl struct {
	pb.UnimplementedStudentService
}

// Hello Hello says hello.
func (s *studentServiceImpl) Hello(
	ctx context.Context,
	req *pb.HelloRequest,
) (*pb.HelloResponse, error) {
	//  增加调用下游服务： 使用下游服务客户端；
	ret, err := call_second_service(ctx, req)
	if err == nil && ret != nil {
		return ret, nil
	}
	rsp := &pb.HelloResponse{
		Msg: "is error response",
	}
	return rsp, err
}
