package main

import (
	pb "grpc-usage-and-gateway/gen/go/proto"
	"log"
	"net"

	"google.golang.org/grpc"
)

func RunGrpcServer() {
	lis, err := net.Listen("tcp", *grpcServerEndpoint)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterGrpcCallDemoServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
