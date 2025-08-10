package main

import (
	"flag"
	"fmt"

	"add_server_on_service_http/internal/config"
	"add_server_on_service_http/internal/handler"
	"add_server_on_service_http/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/foo.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	// 增加自定义配置信息
	fmt.Printf("db config: %+v\n", c.MysqlCfg)
	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
