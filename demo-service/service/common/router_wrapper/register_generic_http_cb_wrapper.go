package router_wrapper

import (
	"github.com/achilsh/go-dev-tools-collection/demo-service/service/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterPostGenericForm[In any, Out any, G gin.IRouter](
	g G,
	url string,
	call func(ctx *gin.Context, inParam *In) (*Out, error),
) {
	g.POST(url, WrapperGenericFormClient(call, middleware.CheckBeforeReadBodyOnForm))
}

func RegisterPostGenericProcess[In any, Out any, G gin.IRouter](
	g G,
	url string,
	call func(ctx *gin.Context, inParam *In) (*Out, error),
) {
	g.POST(url, WrapperGenericClient(call, middleware.CheckBeforeReadBodyOnJson))
}

func RegisterPostGenericNoInProcess[Out any, G gin.IRouter](
	g G,
	url string,
	call func(ctx *gin.Context) (*Out, error),
) {
	g.POST(url, WrapperGenericClientNoIN(call, middleware.CheckBeforeReadBodyOnJson))
}

func RegisterPostGenericNoOutProcess[In any, G gin.IRouter](
	g G,
	url string,
	call func(ctx *gin.Context, inParam *In) error,
) {
	g.POST(url, WrapperGenericClientNoOUT(call, middleware.CheckBeforeReadBodyOnJson))
}

func RegisterGetGenericProcess[In any, Out any, G gin.IRouter](
	g G,
	url string,
	call func(ctx *gin.Context, inParam *In) (*Out, error),
) {
	g.GET(url, WrapperGenericClient(call))
}

func RegisterGetNoInGenericProcess[Out any, G gin.IRouter](g G, url string, call func(ctx *gin.Context) (*Out, error)) {
	g.GET(url, WrapperGenericClientNoIN(call))
}

func RegisterGetNoOutGenericProcess[In any, G gin.IRouter](
	g G,
	url string,
	call func(ctx *gin.Context, inParam *In) error,
) {
	g.GET(url, WrapperGenericClientNoOUT(call))
}
