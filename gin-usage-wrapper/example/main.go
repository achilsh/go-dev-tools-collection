package main

import (
	"log"

	rw "github.com/achilsh/go-dev-tools-collection/gin-usage-wrapper/router_wrapper"

	ginWrapper "github.com/achilsh/go-dev-tools-collection/gin-usage-wrapper/gin_router_wrapper"
	"github.com/achilsh/go-dev-tools-collection/gin-usage-wrapper/middleware"
	"github.com/gin-gonic/gin"
)

type X1V1RequestMsg struct {
	A1 int `json:"a1"`
}
type X1V1Response struct {
	O1 int `json:"o1"`
}

func handlex1_v1(ctx *gin.Context, in *X1V1RequestMsg) (*X1V1Response, error) {
	log.Printf("x1 req: %+v", *in)
	ret := &X1V1Response{
		O1: in.A1 + 1000,
	}
	return ret, nil
}

func Handles(r *ginWrapper.Router) {
	//如果需要在业务逻辑处理之前增加额外的逻辑，只需如下提前注册额外逻辑函数
	rw.SetBeforeReadBodyPostHandlerGetter(
		func() []func(ctx *gin.Context) (int, error) {
			return []func(ctx *gin.Context) (int, error){ /** TODO: add func(ctx* gin.Context)(int, error) **/ }
		},
	)
	//  注册业务逻辑处理函数
	rw.RegisterPostGenericProcess(r.Engine, "/xyz/x1/v1", handlex1_v1)
}

// 请求 http 示例:
// curl -X POST http://127.0.0.1:8080/xyz/x1/v1 -H "Content-Type: application/json" -d '{"a1":123}'
func main() {
	r := ginWrapper.NewRouter(middleware.ParseBody())
	if r == nil {
		log.Fatal("create router fail")
		return
	}
	Handles(r)

	r.Run(8080)
}
