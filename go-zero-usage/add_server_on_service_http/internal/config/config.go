package config

import "github.com/zeromicro/go-zero/rest"

// Mysql:
//
//	IP: 128.1.1.1
//	Port: 13360
//	DB: demo_test_db
type MysqlConf struct {
	Ip     string `json:"IP"`
	Port   int    `json:"Port"`
	DbName string `json:"DB"`
}
type Config struct {
	rest.RestConf
	Auth struct {
		AccessSecret string
		AccessExpire int64
	}
	MysqlCfg MysqlConf `json:"Mysql"`
}
