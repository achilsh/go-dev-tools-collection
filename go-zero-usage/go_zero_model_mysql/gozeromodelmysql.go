package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	"go_zero_model_mysql/go_zero_model_mysql"
	"go_zero_model_mysql/internal/config"
	"go_zero_model_mysql/internal/logic"
	"go_zero_model_mysql/internal/server"
	"go_zero_model_mysql/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/gozeromodelmysql.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		go_zero_model_mysql.RegisterGoZeroModelMysqlServer(grpcServer, server.NewGoZeroModelMysqlServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	//  增加测试
	go func() {
		time.Sleep(2 * time.Second)

		logic.AddInsert(ctx, logx.WithContext(context.Background()))
	}()

	//增加 redis op test
	go func() {
		time.Sleep(2 * time.Second)
		e := ctx.RedisConn.SetCtx(context.Background(), "key_1", "value_1")
		if e != nil {
			fmt.Println("insert redis key fail, err: ", e)
		}
	}()
	s.Start()
}
