package biz_redis

import (
	logger "github.com/achilsh/go-dev-tools-collection/base-lib/log"

	"github.com/achilsh/go-dev-tools-collection/demo-service/dal/redis"
	"github.com/achilsh/go-dev-tools-collection/demo-service/model/user"
)

func GetUserTokenCache(key string) *user.TokenItem {
	items, err := redis.GetRedisInstance().HGetAll(key)
	if err != nil {
		logger.Errorf("get redis hash all items fail, err: %", err)
		return nil
	}

	tokes, err := redis.TransMapToStruct[user.TokenItem](items)
	if err != nil {
		logger.Errorf("trans map to token item fail, err: %v", err)
		return nil
	}
	return tokes
}
