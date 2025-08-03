package main

import (
	"step3/internal/service"
	pb "step3/pb/api/helloworld"

	trpc "trpc.group/trpc-go/trpc-go"
	"trpc.group/trpc-go/trpc-go/log"
)

func main() {
	// init trpc server
	server := trpc.NewServer()
	// Register the greeter service with the server（RESTful HTTP服务）
	pb.RegisterHelloworldService(server, service.NewHelloworldService())
	// Run the server
	if err := server.Serve(); err != nil {
		log.Fatal(err)
	}
}
