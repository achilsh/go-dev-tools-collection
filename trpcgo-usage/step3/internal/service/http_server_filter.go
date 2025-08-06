package service

import (
	"context"
	"time"

	"trpc.group/trpc-go/trpc-go/filter"
	"trpc.group/trpc-go/trpc-go/log"
)

func init() {
	// client.DefaultClient
	// globalService := server.DefaultService()
	// server.RegisterFilter("log-time", ServerFilter)
	// server.RegisterStreamFilter()
	filter.Register("server_cost", ServerFilter, nil) // 拦截器名字自己随便定义，供后续配置文件使用，必须放在 trpc.NewServer() 之前
}

// ServerFilter server 耗时统计，从收到请求到返回响应的处理时间
func ServerFilter(ctx context.Context, req interface{}, next filter.ServerHandleFunc) (rsp interface{}, err error) {

	log.InfoContext(ctx, "server filter start ")
	begin := time.Now() // 业务逻辑处理前打点记录时间戳

	rsp, err = next(ctx, req) // 注意这里必须用户自己调用下一个拦截器，除非有特定目的需要直接返回

	cost := time.Since(begin) // 业务逻辑处理后计算耗时

	log.InfoContext(ctx, "server filter end, cost:  ", cost)
	// 上报耗时到具体监控平台

	return rsp, err // 必须返回 next 的 rsp 和 err，要格外注意不要被自己的逻辑的 rsp 和 err 覆盖

}
