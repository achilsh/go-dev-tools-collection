package main

import (
	"log"

	cliWrap "github.com/achilsh/go-dev-tools-collection/gin-usage-wrapper/client_process_wrapper"
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
	cliWrap.RegisterPostInOut(r.Engine, "/xyz/x1/v1", handlex1_v1)
}
func main() {
	r := ginWrapper.NewRouter(middleware.ParseBody())
	if r == nil {
		log.Fatal("create router fail")
		return
	}
	Handles(r)

	r.Run(8080)
}
