package main

import (
	pb "hello_world_demo/my_hello_world"

	_ "trpc.group/trpc-go/trpc-filter/debuglog"
	_ "trpc.group/trpc-go/trpc-filter/recovery"
	trpc "trpc.group/trpc-go/trpc-go"
	"trpc.group/trpc-go/trpc-go/log"
)

func main() {
	s := trpc.NewServer()
	pb.RegisterStudentServiceService(s.Service("trpc.first.student.StudentService"), &studentServiceImpl{})
	if err := s.Serve(); err != nil {
		log.Fatal(err)
	}
}
