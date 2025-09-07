package bizconfig

import (
	"context"
	"fmt"
	"sync/atomic"

	"gopkg.in/yaml.v3"
	_ "trpc.group/trpc-go/trpc-config-etcd"
	"trpc.group/trpc-go/trpc-go/config"
)

var BizConfigItem *YamlFile = nil

var cfg atomic.Value // 并发安全的 Value

// 定义业务配置内容
type YamlFile struct {
	Custom struct {
		Test    string `yaml:"test"`
		TestObj struct {
			Key1 string `yaml:"key1"`
			Key2 bool   `yaml:"key2"`
			Key3 int32  `yaml:"key3"`
		} `yaml:"test_obj"`
	} `yaml:"custom"`
}

const (
	ConfigTypeValue   = "etcd"
	ConfigBizFileName = "biz_config.yaml"
)

func GetBizConfig() *YamlFile {
	item, ok := cfg.Load().(*YamlFile)
	if !ok {
		return nil
	}
	return item
}
func WatchBizConfig() {
	// 使用 trpc-go/config 中 Watch 接口监听 etcd 远程配置变化
	c, _ := config.Get(ConfigTypeValue).Watch(context.TODO(), ConfigBizFileName)

	go func() {
		for r := range c {
			yf := &YamlFile{}
			fmt.Printf("event: %d, value: %s", r.Event(), r.Value())

			if err := yaml.Unmarshal([]byte(r.Value()), yf); err == nil {
				cfg.Store(yf)
			}
		}
	}()
}
