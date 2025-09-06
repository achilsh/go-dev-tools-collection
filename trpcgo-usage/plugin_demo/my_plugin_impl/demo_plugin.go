package mypluginimpl

import (
	"trpc.group/trpc-go/trpc-go/log"

	"trpc.group/trpc-go/trpc-go/plugin"
)

const (
	pluginName = "my_plugin_demo_1"
	pluginType = "my_plugin_demo_1"
)

var c customPlugin

// 创建一个 func init() 函数
// 通过 Register 注册你的插件。这样别人用你的插件时，
// 只需要在代码中匿名 import 你的包即可
func init() {
	// 该插件调用 plugin.Register 把自己插件注册到 plugin 包。
	plugin.Register(pluginName, &customPlugin{})
}

// customConfig plugin config ；需要配置的插件
type customConfig struct {
	Test    string `yaml:"test"`
	TestObj struct {
		Key1 string `yaml:"key1"`
		Key2 bool   `yaml:"key2"`
		Key3 int32  `yaml:"key3"`
	} `yaml:"test_obj"`
}

// customPlugin struct implements plugin.Factory interface.
type customPlugin struct {
	config customConfig
}

// 对于插件的类型，plugin 包并没有做限制，你可以自行添加插件类型
// Type return plugin type ； 实现了 Factory 接口 的 Type() string 方法
func (custom *customPlugin) Type() string {
	return pluginType
}

// Setup loads plugin by configuration.
// The data structure of the configuration of the plugin needs to be defined in advance。用户自己先定义好具体插件的配置数据结构
// Setup init plugin
// trpc will call Setup function to init plugin.; 实现了 Factory 接口 的 Setup() 方法
func (custom *customPlugin) Setup(name string, decoder plugin.Decoder) error {

	if err := decoder.Decode(&c.config); err != nil {
		return err
	}

	log.Infof("[plugin] init customPlugin success, config: %+v", c.config)

	return nil
}

// Record is a custom plugin function
// you can call this function in your code print plugin config.
func Record() {
	log.Infof("[plugin] call key1 : %s, key2 : %t, key3 : %d",
		c.config.TestObj.Key1, c.config.TestObj.Key2, c.config.TestObj.Key3)
}
