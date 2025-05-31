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


type DemoZRequest struct {
	AA int `json:"aa"`
	Age int `json:"age"`
}
type DemoZResponse struct {
	Address string `json:"address"`
}

var (
	httpClientHandler = httpCli.HttpClientCallTpl[DemoPostReqBody, DemoPostResponseBody]
	httpCliR1Handler  = httpCli.HttpClientCallTpl[DemoPostReqBody, DemoPostRespTplBody[*DemoPostR1Body]]
	httpCliZ1ABCHandler = httpCli.HttpClientCallTpl[DemoZRequest, DemoZResponse]
)

func main() {
	mock_log.LoggerMock()

	in := DemoZRequest{
		AA:   100,
		Age:  100,
	}

	var (
		basUrl    = "http://127.0.0.1:5657"
		urlSuffix = "/demo-server/v1/z1/abc/xxx"
	)

	// response, err := httpCliZ1ABCHandler(
	// 	context.Background(),
	// 	&in,
	// 	httpCli.WithBaseUrl(basUrl),
	// 	httpCli.WithDebug(false),
	// 	httpCli.WithTimeOut(1*time.Second),
	// 	httpCli.WithUrl(urlSuffix),
	// )
	// if err != nil {
	// 	logger.Errorf("http client call fail, err: %v", err)
	// 	return
	// }
	// logger.Debugf("response data: %+v", response)

	r1, err := httpCliZ1ABCHandler(context.Background(), &in, httpCli.WithBaseUrl(basUrl),
		httpCli.WithDebug(false),
		httpCli.WithTimeOut(100*time.Second),
		httpCli.WithUrl(urlSuffix))
	if err != nil {
		logger.Errorf("http client tpl fail, err: %v", err)
		return
	}
	logger.Infof("||||: %+v",*r1)
}
