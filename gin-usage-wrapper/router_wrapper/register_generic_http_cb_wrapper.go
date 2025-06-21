package router_wrapper

import (
	"github.com/gin-gonic/gin"
)

// 如果要修改 OnFormBeforeReadBodyGetter， 在调用下面函数之前，先调用：
//
//	SetBeforeReadBodyOnFormHandlerGetter()
//	SetBeforeReadBodyPostHandlerGetter()
//
// 来修改。
func RegisterPostGenericForm[In any, Out any, G gin.IRouter](
	g G,
	url string,
	call func(ctx *gin.Context, inParam *In) (*Out, error),
) {
	handles := OnFormBeforeReadBodyGetter()
	g.POST(url, WrapperGenericFormClient(call, handles...))
}

func RegisterPostGenericProcess[In any, Out any, G gin.IRouter](
	g G,
	url string,
	call func(ctx *gin.Context, inParam *In) (*Out, error),
) {
	handles := OnPostBeforeReadBodyGetter()
	g.POST(url, WrapperGenericClient(call, handles...))
}

func RegisterPostGenericNoInProcess[Out any, G gin.IRouter](
	g G,
	url string,
	call func(ctx *gin.Context) (*Out, error),
) {
	handles := OnPostBeforeReadBodyGetter()
	g.POST(url, WrapperGenericClientNoIN(call, handles...))
}

func RegisterPostGenericNoOutProcess[In any, G gin.IRouter](
	g G,
	url string,
	call func(ctx *gin.Context, inParam *In) error,
) {
	handles := OnPostBeforeReadBodyGetter()
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
