package main

import (
	pb "student/helloworld"

	_ "trpc.group/trpc-go/trpc-filter/debuglog"
	_ "trpc.group/trpc-go/trpc-filter/recovery"
	trpc "trpc.group/trpc-go/trpc-go"
	"trpc.group/trpc-go/trpc-go/log"

	//  引入 etcd 的服务注册发现插件
	_ "trpc.group/trpc-go/trpc-naming-etcd"
	_ "trpc.group/trpc-go/trpc-naming-etcd/registry"
)

func main() {
	s := trpc.NewServer()
	pb.RegisterHelloWorldServiceService(s, &helloWorldServiceImpl{})
	if err := s.Serve(); err != nil {
		log.Fatal(err)
	}
}
