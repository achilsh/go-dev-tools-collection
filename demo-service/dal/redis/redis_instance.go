package redis

import (
	"fmt"

	"github.com/achilsh/go-dev-tools-collection/demo-service/service/utils/config"
)

var (
	redisInstance RedisOps = nil
)

func InitRedisRes() {
	redisInstance = NewRedisClient(
		WithRedisAddrOpt(fmt.Sprintf("%v:%d", config.GetGlobalConfig().RedisItem.Host, config.GetGlobalConfig().RedisItem.Port)),
		WithRedisDB(config.GetGlobalConfig().RedisItem.DataBase), WithRedisPasswd(config.GetGlobalConfig().RedisItem.Password))
}
func GetRedisInstance() RedisOps {
	return redisInstance
}
