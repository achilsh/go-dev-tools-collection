package main

import (
	"context"

	pb "second_demo/helloworld"
)

type greeterImpl struct {
	pb.UnimplementedGreeter
}

func (s *greeterImpl) SayHello(
	ctx context.Context,
	req *pb.HelloRequest,
) (*pb.HelloReply, error) {
	rsp := &pb.HelloReply{
		Msg: "this is greeter msg response",
	}
	return rsp, nil
}
