package main

import (
	"flag"

	"demo-service/dal/redis"
	"demo-service/router"
	"demo-service/service/utils/config"
	"demo-service/service/utils/init_res"
)

var aiConfigFile *string

func init() {
	aiConfigFile = flag.String("f", "../../conf/demo-server.yaml", "demo server config file, absolute path to config file")

	flag.Parse()
}

func main() {
	if err := init_res.InitConfigLog(aiConfigFile); err != nil {
		return
	}
	redis.InitRedisRes()
	//
	router.RegisterRouter()

	router.GetRouter().Run(config.GetGlobalConfig().Server.Port)
}
