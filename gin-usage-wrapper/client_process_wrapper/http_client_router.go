package clientprocesswrapper

import (
	"github.com/gin-gonic/gin"
)

func RegisterPostInOut[IN any, OUT any, G gin.IRouter](g G, url string, f func(*gin.Context, IN) (OUT, error)) {
	g.POST(url, ClientProcessWrapper(f))
}

func RegisterPostIn[IN any, G gin.IRouter](g G, url string, f func(*gin.Context, IN) error) {
	g.POST(url, ClientProcessWrapper(f))
}

func RegisterPostOut[OUT any, G gin.IRouter](g G, url string, f func(*gin.Context) (OUT, error)) {
	g.POST(url, ClientProcessWrapper(f))
}

func RegisterGetInOut[IN any, OUT any, G gin.IRouter](g G, url string, f func(*gin.Context, IN) (OUT, error)) {
	g.GET(url, ClientProcessWrapper(f))
}

func RegisterGetIn[IN any, G gin.IRouter](g G, url string, f func(*gin.Context, IN) error) {
	g.GET(url, ClientProcessWrapper(f))
}

func RegisterGetOut[OUT any, G gin.IRouter](g G, url string, f func(*gin.Context) (OUT, error)) {
	g.GET(url, ClientProcessWrapper(f))
}
