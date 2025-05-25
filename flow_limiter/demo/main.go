package main

import (
	"time"

	log "github.com/achilsh/go-dev-tools-collection/base-lib/log"
	logger "github.com/achilsh/go-dev-tools-collection/demo-service/service/utils/mock/mock_log"
	tl "github.com/achilsh/go-dev-tools-collection/flow_limiter/token_limit"
	rdslib "github.com/achilsh/go-dev-tools-collection/redis-wrapper/lib"
)

var limiterHandle *tl.TokenLimiter = nil

//
func InitLimiter() {
	limiterHandle = tl.NewTokenLimiter(
	tl.InitRedis(rdslib.WithRedisAddrOpt("0.0.0.0:6379")), tl.LoadHashLua())
}

func main() {
	logger.LoggerMock()
	InitLimiter()

	for i := 0; i < 1000; i++ {
		// 每分钟最多 60 个请求（即 maxTokens=60, refillRate=1）
		numMaxPerSecond := 5  //每秒消耗token最大数
		fillPerSecond :=1 //每秒填充token数
		ok := limiterHandle.Allow(tl.InitKeys("a2_request", "key_total_nums"),tl.InitControlInfo(numMaxPerSecond, fillPerSecond, 1000))
		time.Sleep(100 *time.Millisecond)
		if !ok {
			log.Infof("not allow. %+v", time.Now())
		} else {
			log.Infof("is allow: %+v", time.Now())
		}
	}
}