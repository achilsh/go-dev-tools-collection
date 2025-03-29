package ginrouterwrapper

import (
	"fmt"
	"log"
	"time"

	"github.com/achilsh/http-graceful/graceful.v1"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

type Router struct {
	*gin.Engine
}

// NewRouter 创建一个router;封装了 gin.Engine. preProcess 可以类似 recovery, cors等处理
func NewRouter(preProcess ...gin.HandlerFunc) *Router {
	gin.SetMode(gin.DebugMode)
	router := gin.New()

	router.ContextWithFallback = true
	pprof.Register(router)
	if len(preProcess) > 0 {
		router.Use(preProcess...)
	}

	return &Router{Engine: router}
}

// Run 启动gin 服务.
func (r *Router) Run(port int) {
	log.Println("Router starting")
	graceful.Run(fmt.Sprintf(":%v", port), 30*time.Second, r.Engine)
}
