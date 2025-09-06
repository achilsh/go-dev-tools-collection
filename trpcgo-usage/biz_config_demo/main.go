package main

import (
	"fmt"
	pb "hello_world_demo/my_hello_world"

	_ "trpc.group/trpc-go/trpc-filter/debuglog"
	_ "trpc.group/trpc-go/trpc-filter/recovery"
	trpc "trpc.group/trpc-go/trpc-go"
	"trpc.group/trpc-go/trpc-go/config"
	"trpc.group/trpc-go/trpc-go/log"
)

// customStruct it defines the struct of the custom configuration file read.
type customStruct struct {
	Custom struct {
		Test    string `yaml:"test"`
		TestObj struct {
			Key1 string `yaml:"key1"`
			Key2 bool   `yaml:"key2"`
			Key3 int32  `yaml:"key3"`
		} `yaml:"test_obj"`
	} `yaml:"custom"`
}

//  load buzi config

func loadBizConfig() {
	//  codec 支持 yaml, json, toml
	// provider 支持 file, etcd 等。
	c, err := config.Load("custom.yaml", config.WithCodec("yaml"), config.WithProvider("file"))
	if err != nil {
		log.Fatal(fmt.Sprintf("load custom.yaml file fail: %v", err))
		return
	}
	//print buzi config item :
	ctest := c.GetString("custom.test", "")
	log.Debugf("!!-----custom.test item value: %v", ctest)

	//把业务配置文件反序列化为内存结构体，eg：
	var confItems customStruct
	if err := c.Unmarshal(&confItems); err != nil {
		log.Errorf("unmarshal buzi config to struct fail, err: %v", err)
		return
	}

	log.Debugf("2------> config content: %+v", confItems)

}
func main() {

	s := trpc.NewServer()

	loadBizConfig()

	pb.RegisterHelloWorldServiceService(s.Service("my_hello_world.HelloWorldService"), &helloWorldServiceImpl{})
	if err := s.Serve(); err != nil {
		log.Fatal(err)
	}
}
