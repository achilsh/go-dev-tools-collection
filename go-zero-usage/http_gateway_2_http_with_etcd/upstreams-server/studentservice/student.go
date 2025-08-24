package main

import (
	"flag"
	"fmt"

	"http_gateway_demo/upstreams-server/studentservice/internal/config"
	"http_gateway_demo/upstreams-server/studentservice/internal/handler"
	"http_gateway_demo/upstreams-server/studentservice/internal/svc"

	"github.com/zeromicro/zero-contrib/rest/registry/etcd"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/student-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)

	// http 自身向 etcd 注册动作
	logx.Must(etcd.RegisterRest(c.Etcd, c.RestConf))

	server.Start()
}
