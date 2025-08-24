package main

import (
	"flag"
	"fmt"

	"http_gateway_demo/upstreams-server/userservice/internal/config"
	"http_gateway_demo/upstreams-server/userservice/internal/handler"
	"http_gateway_demo/upstreams-server/userservice/internal/middleware"
	"http_gateway_demo/upstreams-server/userservice/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/userservice-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	server.Use(middleware.HeaderMiddlewares)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)

	server.Start()
}
