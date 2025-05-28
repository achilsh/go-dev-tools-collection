package main

import (
	"context"
	"time"

	logger "github.com/achilsh/go-dev-tools-collection/base-lib/log"
	"github.com/achilsh/go-dev-tools-collection/demo-service/service/utils/mock/mock_log"
	httpCli "github.com/achilsh/go-dev-tools-collection/http_client_wrapper"
)

// 定义请求和回包数据结构

type DemoPostReqBody struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type DemoPostResponseBody struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Award   int    `json:"award"`
}

// 返回数据也是多个类型
type DemoPostRespTplBody[T any] struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    []T    `json:"data"`
}

type DemoPostR1Body struct {
	Award   int    `json:"award"`
	Address string `json:"address"`
}

var (
	httpClientHandler = httpCli.HttpClientCallTpl[DemoPostReqBody, DemoPostResponseBody]
	httpCliR1Handler  = httpCli.HttpClientCallTpl[DemoPostReqBody, DemoPostRespTplBody[*DemoPostR1Body]]
)

func main() {
	mock_log.LoggerMock()

	in := DemoPostReqBody{
		Id:   100,
		Name: "test_demo_call",
		Age:  100,
	}

	var (
		basUrl    = "http://xxx.com"
		urlSuffix = "/api/v1/demo"
	)

	response, err := httpClientHandler(
		context.Background(),
		&in,
		httpCli.WithBaseUrl(basUrl),
		httpCli.WithDebug(true),
		httpCli.WithTimeOut(1*time.Second),
		httpCli.WithUrl(urlSuffix),
	)
	if err != nil {
		logger.Errorf("http client call fail, err: %v", err)
		return
	}
	logger.Debugf("response data: %+v", response)

	r1, err := httpCliR1Handler(context.Background(), &in, httpCli.WithBaseUrl(basUrl),
		httpCli.WithDebug(true),
		httpCli.WithTimeOut(1*time.Second),
		httpCli.WithUrl(urlSuffix))
	if err != nil {
		logger.Errorf("http client tpl fail, err: %v", err)
		return
	}
	_ = r1
}
