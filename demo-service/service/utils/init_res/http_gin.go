package init_res

import (
	"fmt"
	"time"

	"github.com/achilsh/http-graceful/graceful.v1"
	"github.com/gin-gonic/gin"

	logger "github.com/achilsh/go-dev-tools-collection/base-lib/log"
)

type Router struct {
	*gin.Engine
}

func (r *Router) GetGin() *gin.Engine {
	return r.Engine
}

func (r *Router) Run(port int) {
	logger.Debugf("http  server start begin! port: %v", port)
	graceful.Run(fmt.Sprintf(":%v", port), 30*time.Second, r.Engine)
}
