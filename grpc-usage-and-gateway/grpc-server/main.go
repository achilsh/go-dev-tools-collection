package main

import (
	"flag"

	"google.golang.org/grpc/grpclog"
)

func main() {
	flag.Parse()

	go func() {
		RunGrpcServer()
	}()

	if err := RunServer(); err != nil {
		grpclog.Fatal(err)
	}
}
