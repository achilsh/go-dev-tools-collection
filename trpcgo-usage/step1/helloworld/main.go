package main

import (
	pb "step1/helloworld"

	_ "trpc.group/trpc-go/trpc-filter/debuglog"
	_ "trpc.group/trpc-go/trpc-filter/recovery"
	trpc "trpc.group/trpc-go/trpc-go"
	"trpc.group/trpc-go/trpc-go/log"
)

func main() {
	// 解析 trpc_go.yaml文件内容，启动服务器，从内容中启动多个服务，可以使用 -conf 指定配置文件。
	s := trpc.NewServer() //里面可选参数有些什么？

	// 根据服务名 helloworld.HelloWorldService 返回一个服务器，把具体实现的服务注册到服务器中。
	// 将 trpc 服务和具体业务实现绑定在了一起，启动服务就对接上了业务逻辑; s.Service("helloworld.HelloWorldService")
	pb.RegisterHelloWorldServiceService(s, &helloWorldServiceImpl{})

	// 启动注册到该服务器中的服务
	if err := s.Serve(); err != nil {
		log.Fatal(err)
	}
}
