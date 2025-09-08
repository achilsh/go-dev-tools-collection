package main

import (
	"context"

	pb "helloworld/pb"
)

type greeterImpl struct {
	pb.UnimplementedGreeter
}

func (s *greeterImpl) Hello(
	ctx context.Context,
	req *pb.HelloRequest,
) (*pb.HelloReply, error) {
	rsp := &pb.HelloReply{
		Msg: "this is default data.",
	}
	return rsp, nil
}
