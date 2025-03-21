package file

import (
	"errors"

	conf "github.com/achilsh/go-dev-tools-collection/base-lib/config"
	"github.com/spf13/viper"
	"go-micro.dev/v4/config"
	"go-micro.dev/v4/config/reader"
	"go-micro.dev/v4/logger"
)

type File struct {
	Path string
}

// Init 初始化配置
//
// cfg: 配置结构体
func (p *File) Init(cfg interface{}) error {
	// 初始化配置
	viper.SetConfigFile(p.Path)
	err := viper.ReadInConfig()
	if err != nil {
		return errors.New("read config file error: " + err.Error())
	}
	err = viper.Unmarshal(cfg)
	if err != nil {
		return errors.New("unmarshal config file error: " + err.Error())
	}
	logger.Debug(cfg)
	return nil
}

// InitWithWatch 初始化配置并监听配置变化
//
// File无法监听配置变化, callback无效
func (p *File) InitWithWatch(cfg interface{}, callback func(value reader.Value)) (config.Watcher, error) {
	err := p.Init(cfg)
	return nil, err
}

// Write 更新并写入配置
func (p *File) Write(key string, value interface{}) error {
	if p.Path != "" {
		viper.SetConfigFile(p.Path)
	}
	viper.Set(key, value)
	return viper.WriteConfig()
}

// NewConfig 从file读取配置
//
// path: 配置所在路径
func NewConfig(path string) conf.Config {
	return &File{
		Path: path,
	}
}
