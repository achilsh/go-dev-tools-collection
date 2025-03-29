package router

import (
	"net/http"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"

	"demo-service/service/middleware"
	"demo-service/service/utils/init_res"
)

var globalRouter *init_res.Router = nil

func init() {
	GetRouter()
}

func GetRouter() *init_res.Router {
	if globalRouter == nil {
		globalRouter = newRouter()
	}
	return globalRouter
}

// newRouter 创建 gin 的对象
func newRouter() *init_res.Router {
	router := gin.New()
	router.ContextWithFallback = true
	pprof.Register(router)

	router.Use(
		middleware.HandlerW2(),
		middleware.HandleOther(),
		middleware.LanguageManger(),
		middleware.CORS(),
		middleware.AccessLogMiddleware(),
		middleware.Recovery())

	router.GET("/", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "ok")
	})
	router.GET("/test", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "ok")
	})

	return &init_res.Router{
		Engine: router,
	}
}
