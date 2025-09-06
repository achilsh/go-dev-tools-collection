package main

import (
	pb "hello_world_demo/my_hello_world"

	_ "plugin_demo/my_plugin_impl"
	mypluginimpl "plugin_demo/my_plugin_impl"

	_ "trpc.group/trpc-go/trpc-filter/debuglog"
	_ "trpc.group/trpc-go/trpc-filter/recovery"
	trpc "trpc.group/trpc-go/trpc-go"
	"trpc.group/trpc-go/trpc-go/log"
)

func main() {
	s := trpc.NewServer()
	pb.RegisterHelloWorldServiceService(s.Service("my_hello_world.HelloWorldService"), &helloWorldServiceImpl{})

	mypluginimpl.Record()
	if err := s.Serve(); err != nil {
		log.Fatal(err)
	}
}
