package main

import (
	"context"

	pb "hello_world_demo/my_hello_world"

	"trpc.group/trpc-go/tnet/log"
)

type helloWorldServiceImpl struct {
	pb.UnimplementedHelloWorldService
}

// Hello Hello says hello.
func (s *helloWorldServiceImpl) Hello(
	ctx context.Context,
	req *pb.HelloRequest,
) (*pb.HelloResponse, error) {
	log.Infof("req: %+v", *req)
	rsp := &pb.HelloResponse{}
	return rsp, nil
}
