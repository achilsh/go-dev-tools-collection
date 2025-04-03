package main

import (
	"time"

	lib "github.com/achilsh/go-dev-tools-collection/redis-wrapper/lib"
)

func main() {
	cli := lib.NewRedisClient(lib.WithRedisAddrOpt("10.240.34.36:30356"),
		lib.WithRedisPasswd("xxxx"), lib.WithRedisDB(15))
	cli.SetEx("aaa", "aaa", 10*time.Second)

}
