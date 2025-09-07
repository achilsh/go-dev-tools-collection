package main

import (
	"context"

	pb "second_demo/helloworld"
)

type helloImpl struct {
	pb.UnimplementedHello
}

func (s *helloImpl) SayHi(
	ctx context.Context,
	req *pb.HelloRequest,
) (*pb.HelloReply, error) {
	rsp := &pb.HelloReply{
		Msg: "this is hello response.",
	}
	return rsp, nil
}
