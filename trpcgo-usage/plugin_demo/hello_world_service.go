package main

import (
	"context"

	pb "hello_world_demo/my_hello_world"
)

type helloWorldServiceImpl struct {
	pb.UnimplementedHelloWorldService
}

// Hello Hello says hello.
func (s *helloWorldServiceImpl) Hello(
	ctx context.Context,
	req *pb.HelloRequest,
) (*pb.HelloResponse, error) {
	rsp := &pb.HelloResponse{}
	return rsp, nil
}
