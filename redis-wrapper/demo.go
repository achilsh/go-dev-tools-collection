package main

import (
	"time"

	"redis-wrapper-demo/lib"
)

func main() {
	cli := lib.NewRedisClient(lib.WithRedisAddrOpt("10.240.34.36:30356"),
		lib.WithRedisPasswd("xxxx"), lib.WithRedisDB(15))
	cli.SetEx("aaa", "aaa", 10*time.Second)

}
