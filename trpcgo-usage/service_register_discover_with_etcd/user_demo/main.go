package main

import (
	"context"
	"time"
	pb "user/helloworld"

	_ "trpc.group/trpc-go/trpc-filter/debuglog"
	_ "trpc.group/trpc-go/trpc-filter/recovery"
	"trpc.group/trpc-go/trpc-go"
	"trpc.group/trpc-go/trpc-go/log"

	//  集成 etcd 服务发现的插件
	pbStud "student/helloworld"

	_ "trpc.group/trpc-go/trpc-naming-etcd"
	_ "trpc.group/trpc-go/trpc-naming-etcd/registry"
)

func mockCallUpstreamByEtcdDiscover() {
	go func() {
		for i := 0; i < 10; i++ {
			time.Sleep(1 * time.Second)

			clientProxy := pbStud.NewHelloWorldServiceClientProxy()
			req := &pbStud.HelloRequest{
				Msg: "hello",
			}

			rsp, err := clientProxy.Hello(context.Background(), req)
			if err != nil {
				log.Error(err.Error())
				return
			}

			log.Info("req:%v, rsp:%v, err:%v", req, rsp, err)

		}

	}()

}
func main() {
	s := trpc.NewServer()
	mockCallUpstreamByEtcdDiscover()

	pb.RegisterHelloWorldServiceService(s, &helloWorldServiceImpl{})
	if err := s.Serve(); err != nil {
		log.Fatal(err)
	}
}
