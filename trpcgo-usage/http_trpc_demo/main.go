package main

import (
	pb "helloworld/pb"

	_ "trpc.group/trpc-go/trpc-filter/debuglog"
	_ "trpc.group/trpc-go/trpc-filter/recovery"
	trpc "trpc.group/trpc-go/trpc-go"
	"trpc.group/trpc-go/trpc-go/log"
)

func main() {
	s := trpc.NewServer()
	pb.RegisterGreeterService(s.Service("trpc.helloworld.Greeter"), &greeterImpl{})
	pb.RegisterGreeterService(s.Service("http.helloworld.Greeter"), &greeterImpl{})
	// pb.RegisterGreeterService(s, &greeterImpl{})
	if err := s.Serve(); err != nil {
		log.Fatal(err)
	}
}
