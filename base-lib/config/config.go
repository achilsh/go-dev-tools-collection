package config

import (
	"go-micro.dev/v4/config"
	"go-micro.dev/v4/config/reader"
)

type Config interface {
	// Init 初始化配置
	Init(cfg interface{}) error
	// InitWithWatch 初始化配置并监听配置变化
	InitWithWatch(cfg interface{}, callback func(value reader.Value)) (config.Watcher, error)
	// Write 更新并写入配置
	Write(key string, value interface{}) error
}
