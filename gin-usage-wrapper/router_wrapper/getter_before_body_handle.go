package router_wrapper

import (
	"github.com/gin-gonic/gin"
)

var OnFormBeforeReadBodyGetter = func() []func(ctx *gin.Context) (int, error) {
	return []func(ctx *gin.Context) (int, error){
		// middleware.CheckBeforeReadBodyOnForm,
	}
}
var OnPostBeforeReadBodyGetter = func() []func(ctx *gin.Context) (int, error) {
	return []func(ctx *gin.Context) (int, error){
		// middleware.CheckBeforeReadBodyOnJson,
	}
}

// 上面是默认的处理函数。
type BeforeReadBodyHandlersGetter = func() []func(ctx *gin.Context) (int, error)

// 这是供业务方注册处理函数（修改默认值）
func SetBeforeReadBodyOnFormHandlerGetter(getter BeforeReadBodyHandlersGetter) {
	OnFormBeforeReadBodyGetter = getter
}

// 这是供业务方注册处理函数（修改默认值）
func SetBeforeReadBodyPostHandlerGetter(getter BeforeReadBodyHandlersGetter) {
	OnPostBeforeReadBodyGetter = getter
}

func FlagNoOutputLog(ctx *gin.Context) {
	if ctx == nil {
		return
	}
	ctx.Set("omit_log", 1)
}

func IsNoOutputLog(ctx *gin.Context) bool {
	if ctx == nil {
		return false
	}
	_, exist := ctx.Get("omit_log")
	return exist
}
