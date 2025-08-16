package main

import (
	"flag"
	"fmt"

	"http_server_etcd/internal/config"
	"http_server_etcd/internal/handler"
	"http_server_etcd/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/zero-contrib/rest/registry/etcd"
)

var configFile = flag.String("f", "etc/httpserveretcd-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)

	logx.Must(etcd.RegisterRest(c.RegisterEtcd, c.RestConf))

	server.Start()
}
