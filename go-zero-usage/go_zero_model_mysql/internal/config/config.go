package config

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	//  增加数据库连接地址配置：
	DataSource string
	// 增加 redis 配置
	RedisConf redis.RedisConf
}
