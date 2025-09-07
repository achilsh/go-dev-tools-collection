package main

import (
	"context"
	pb "hello_world_demo/my_hello_world"

	secondPb "second_demo/helloworld"

	"trpc.group/trpc-go/tnet/log"
	trpc "trpc.group/trpc-go/trpc-go"
)

func callGreeterSayHello() *secondPb.HelloReply {
	proxy := secondPb.NewGreeterClientProxy(
	// client.WithTarget("ip://127.0.0.1:8000"),
	// client.WithProtocol("trpc"),
	)
	ctx := trpc.BackgroundContext()
	// Example usage of unary client.
	reply, err := proxy.SayHello(ctx, &secondPb.HelloRequest{})
	if err != nil {
		log.Fatalf("err: %v", err)
	}
	log.Debugf("simple  rpc   receive: %+v", reply)
	return reply
}

func callHelloSayHi() *secondPb.HelloReply {
	proxy := secondPb.NewHelloClientProxy(
	// client.WithTarget("ip://127.0.0.1:8001"),
	// client.WithProtocol("trpc"),
	)
	ctx := trpc.BackgroundContext()
	// Example usage of unary client.
	reply, err := proxy.SayHi(ctx, &secondPb.HelloRequest{})
	if err != nil {
		log.Fatalf("err: %v", err)
	}
	log.Debugf("simple  rpc   receive: %+v", reply)
	return reply
}

func call_second_service(ctx context.Context,
	req *pb.HelloRequest,
) (*pb.HelloResponse, error) {

	r1 := callHelloSayHi()
	r2 := callGreeterSayHello()
	var retmsg = ""
	if r1 != nil {
		if len(retmsg) > 0 {
			retmsg += ", " + r1.GetMsg()
		} else {
			retmsg += r1.GetMsg()
		}
	}
	if r2 != nil {
		if len(retmsg) > 0 {
			retmsg += ", " + r2.GetMsg()
		} else {
			retmsg += r2.GetMsg()
		}
	}
	return &pb.HelloResponse{
		Msg: retmsg,
	}, nil
}
