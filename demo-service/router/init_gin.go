package router

import (
	"net/http"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"

	"github.com/achilsh/go-dev-tools-collection/demo-service/service/middleware"
	"github.com/achilsh/go-dev-tools-collection/demo-service/service/utils/init_res"
)

var globalRouter *init_res.Router = nil

type ginConfType struct {
	uploadFileMaxSizeLimit int64 // byte
}

var ginDefaultConf ginConfType

func init() {
	ginDefaultConf = ginConfType{
		uploadFileMaxSizeLimit: 10 * 1024 * 1024,
	}
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
	router.MaxMultipartMemory = ginDefaultConf.uploadFileMaxSizeLimit // 设置文件上传大小的限制 10MB

	return &init_res.Router{
		Engine: router,
	}
}
