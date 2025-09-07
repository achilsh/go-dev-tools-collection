package main

import (
	pb "second_demo/helloworld"

	_ "trpc.group/trpc-go/trpc-filter/debuglog"
	_ "trpc.group/trpc-go/trpc-filter/recovery"
	trpc "trpc.group/trpc-go/trpc-go"
	"trpc.group/trpc-go/trpc-go/log"
)

func main() {
	s := trpc.NewServer()
	//用户需要明确指定 naming service 和 proto service 的映射关系
	pb.RegisterGreeterService(s.Service("trpc.second.helloworld.Greeter"), &greeterImpl{})
	// 把 Proto Service 注册到 Naming Service，多服务场景需要指定 Naming Service 的名称
	pb.RegisterHelloService(s.Service("trpc.second.helloworld.Hello"), &helloImpl{})
	if err := s.Serve(); err != nil {
		log.Fatal(err)
	}
}
