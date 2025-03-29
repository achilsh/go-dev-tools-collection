package router_wrapper

import (
	"github.com/gin-gonic/gin"
)

// RegisterPostProcess post 请求， 有具体的业务请求体和业务回包体
func RegisterPostProcess[In any, Out any, G gin.IRouter](g G, url string, call func(ctx *gin.Context, inParam In) (Out, error)) {
	g.POST(url, WrapperClient(call))
}

// RegisterPostNoInProcess post 请求， 没有具体的业务请求体，但有业务回包体
func RegisterPostNoInProcess[Out any, G gin.IRouter](g G, url string, call func(ctx *gin.Context) (Out, error)) {
	g.POST(url, WrapperClient(call))
}

// RegisterPostNoOutProcess post 请求， 有具体的业务请求体， 但是没有业务 回包体
func RegisterPostNoOutProcess[In any, G gin.IRouter](g G, url string, call func(ctx *gin.Context, inParam In) error) {
	g.POST(url, WrapperClient(call))
}

// RegisterGetProcess get 请求，有具体的业务请求体和业务回包体
func RegisterGetProcess[In any, Out any, G gin.IRouter](g G, url string, call func(ctx *gin.Context, inParam In) (Out, error)) {
	g.GET(url, WrapperClient(call))
}

// RegisterGetNoInProcess get 请求，没有具体的业务请求体，但是有业务回包体
func RegisterGetNoInProcess[Out any, G gin.IRouter](g G, url string, call func(ctx *gin.Context) (Out, error)) {
	g.GET(url, WrapperClient(call))
}

// RegisterGetNoOutProcess get 请求， 有具体业务请求体，但是没有业务回包体
func RegisterGetNoOutProcess[In any, G gin.IRouter](g G, url string, call func(ctx *gin.Context, inParam In) error) {
	g.GET(url, WrapperClient(call))
}
