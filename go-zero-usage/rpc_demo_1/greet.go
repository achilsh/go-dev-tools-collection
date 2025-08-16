package main

import (
	"flag"
	"fmt"

	"rpc_demo_1/internal/config"
	greetServer "rpc_demo_1/internal/server/greet"
	messageServer "rpc_demo_1/internal/server/message"
	"rpc_demo_1/internal/svc"
	"rpc_demo_1/pb/gen/greet"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/greet.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		greet.RegisterGreetServer(grpcServer, greetServer.NewGreetServer(ctx))
		greet.RegisterMessageServer(grpcServer, messageServer.NewMessageServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
