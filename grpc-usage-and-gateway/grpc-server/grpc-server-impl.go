package main

import (
	"context"
	pb "grpc-usage-and-gateway/gen/go/proto"
	"log"
)

type server struct {
	pb.UnimplementedGrpcCallDemoServer
}

func (s *server) Echo(_ context.Context, in *pb.StringMessage) (*pb.StringMessage, error) {
	log.Printf("Received: %v", in.GetValue())
	return &pb.StringMessage{Value: "Hello " + in.GetValue()}, nil
}
func (s *server) Hello(_ context.Context, in *pb.StringMessage) (*pb.StringMessage, error) {
	log.Printf("Received hello request.: %v", in.GetValue())
	return &pb.StringMessage{Value: "world " + in.GetValue()}, nil
}
