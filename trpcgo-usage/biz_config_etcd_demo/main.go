package main

import (
	bizconfig "biz_config_etcd_demo/biz_config"
	pb "hello_world_demo/my_hello_world"
	"time"

	_ "trpc.group/trpc-go/trpc-filter/debuglog"
	_ "trpc.group/trpc-go/trpc-filter/recovery"
	trpc "trpc.group/trpc-go/trpc-go"
	"trpc.group/trpc-go/trpc-go/log"
)

func getBizConfig() {
	for {
		time.Sleep(5 * time.Second)
		bizCfgItem := bizconfig.GetBizConfig()
		if bizCfgItem == nil {
			log.Errorf("get biz config content is empty.")
			continue
		}
		log.Debugf("biz config value: %+v", *bizCfgItem)
	}
}

func main() {
	s := trpc.NewServer()
	bizconfig.WatchBizConfig()

	getBizConfig()

	pb.RegisterHelloWorldServiceService(s.Service("my_hello_world.HelloWorldService"), &helloWorldServiceImpl{})
	if err := s.Serve(); err != nil {
		log.Fatal(err)
	}
}
