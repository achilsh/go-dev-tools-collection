package main

import (
	"flag"
	"http_gateway_demo/middleware"
	"http_gateway_demo/rewrites"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/gateway"
)

var configFile = flag.String("f", "etc/gateway.yaml", "config file")

func main() {
	flag.Parse()

	// var c gateway.GatewayConf
	var c rewrites.ServiceConfig
	conf.MustLoad(*configFile, &c)
	gw := gateway.MustNewServer(c.GatewayConf)
	defer gw.Stop()
	gw.Use(middleware.HeaderMiddlewares)
	//
	gw.Use(middleware.PathRewriteMiddleware(rewrites.TransRewriteCfgToRules(&c.RewritHttpUrlPathConf)))
	gw.Start()
}
