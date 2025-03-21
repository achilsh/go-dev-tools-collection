package consul

import (
	"errors"

	conf "github.com/achilsh/go-dev-tools-collection/base-lib/config"
	"github.com/go-micro/plugins/v4/config/encoder/yaml"
	"github.com/go-micro/plugins/v4/config/source/consul"
	"go-micro.dev/v4/config"
	"go-micro.dev/v4/config/reader"
	"go-micro.dev/v4/config/reader/json"
	"go-micro.dev/v4/config/source"
	"go-micro.dev/v4/logger"
)

const (
	DefaultAddr   = "127.0.0.1:8500"
	DefaultPrefix = "/cuav"
	DefaultPath   = "dev"
)

type Consul struct {
	Addr   string
	Prefix string
	Path   string
}

// Init 初始化配置
//
// cfg: 配置结构体
func (p *Consul) Init(cfg interface{}) error {
	// 连接consul
	c, _ := config.NewConfig(config.WithReader(json.NewReader(reader.WithEncoder(yaml.NewEncoder()))))
	defer c.Close()
	s := consul.NewSource(consul.WithAddress(p.Addr), consul.WithPrefix(p.Prefix), consul.StripPrefix(true), source.WithEncoder(yaml.NewEncoder()))
	err := c.Load(s)
	if err != nil {
		return errors.New("connect consul error: " + err.Error())
	}
	// 初始化配置
	err = c.Get(p.Path).Scan(&cfg)
	if err != nil {
		return errors.New("get config from consul error: " + err.Error())
	}
	logger.Debug(cfg)
	return nil
}

// InitWithWatch 初始化配置并监听配置变化
//
// cfg: 配置结构体; callback: 配置变化时的回调函数
func (p *Consul) InitWithWatch(cfg interface{}, callback func(value reader.Value)) (config.Watcher, error) {
	// 连接consul
	c, _ := config.NewConfig(config.WithReader(json.NewReader(reader.WithEncoder(yaml.NewEncoder()))))
	defer c.Close()
	s := consul.NewSource(consul.WithAddress(p.Addr), consul.WithPrefix(p.Prefix), consul.StripPrefix(true), source.WithEncoder(yaml.NewEncoder()))
	err := c.Load(s)
	if err != nil {
		return nil, errors.New("connect consul error: " + err.Error())
	}
	// 初始化配置
	err = c.Get(p.Path).Scan(&cfg)
	if err != nil {
		return nil, errors.New("get config from consul error: " + err.Error())
	}
	logger.Debug(cfg)
	// 监听配置变化
	watch, err := c.Watch(p.Path)
	if err != nil {
		return nil, errors.New("watch config change error: " + err.Error())
	}
	go func() {
		for {
			v, err := watch.Next()
			if err != nil {
				logger.Error("process config change error: ", err)
				watch.Stop()
				return
			}
			// 同步变化
			callback(v)
		}
	}()
	return watch, nil
}

// Write 更新并写入配置
//
// 待实现
func (p *Consul) Write(key string, value interface{}) error {
	return nil
}

// NewConfig 从consul读取配置
//
// addr: consul地址; prefix: 配置路径前缀; path: 配置所在路径
func NewConfig(addr string, prefix string, path string) conf.Config {
	if addr == "" {
		addr = DefaultAddr
	}
	if prefix == "" {
		prefix = DefaultPrefix
	}
	if path == "" {
		path = DefaultPath
	}
	return &Consul{
		Addr:   addr,
		Prefix: prefix,
		Path:   path,
	}
}
