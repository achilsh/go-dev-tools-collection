package router_wrapper

import (
	"github.com/gin-gonic/gin"
)

func RegisterPostGenericForm[In any, Out any, G gin.IRouter](
	g G,
	url string,
	call func(ctx *gin.Context, inParam *In) (*Out, error),
) {
	handles := defaultOnFormBeforeReadBodyGetter()
	g.POST(url, WrapperGenericFormClient(call, handles...))
}

func RegisterPostGenericProcess[In any, Out any, G gin.IRouter](
	g G,
	url string,
	call func(ctx *gin.Context, inParam *In) (*Out, error),
) {
	handles := defaultOnPostBeforeReadBodyGetter()
	g.POST(url, WrapperGenericClient(call, handles...))
}

func RegisterPostGenericNoInProcess[Out any, G gin.IRouter](
	g G,
	url string,
	call func(ctx *gin.Context) (*Out, error),
) {
	handles := defaultOnPostBeforeReadBodyGetter()
	g.POST(url, WrapperGenericClientNoIN(call, handles...))
}

func RegisterPostGenericNoOutProcess[In any, G gin.IRouter](
	g G,
	url string,
	call func(ctx *gin.Context, inParam *In) error,
) {
	handles := defaultOnPostBeforeReadBodyGetter()
	g.POST(url, WrapperGenericClientNoOUT(call, handles...))
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
