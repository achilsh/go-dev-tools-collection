package router_wrapper

import (
	"github.com/gin-gonic/gin"
)

type BeforeReadBodyHandlersGetter  = func() []func(ctx *gin.Context) (int, error)
var defaultOnFormBeforeReadBodyGetter = func() []func(ctx *gin.Context) (int, error) {
	return []func(ctx *gin.Context) (int, error){
		  // middleware.CheckBeforeReadBodyOnForm,
	}
}
var defaultOnPostBeforeReadBodyGetter = func() []func(ctx *gin.Context) (int, error) {
	return []func(ctx *gin.Context)(int, error) {
		// middleware.CheckBeforeReadBodyOnJson,
	}
}

func SetBeforeReadBodyOnFormHandlerGetter(getter BeforeReadBodyHandlersGetter) {
	defaultOnFormBeforeReadBodyGetter = getter
}
func SetBeforeReadBodyPostHandlerGetter(getter BeforeReadBodyHandlersGetter) {
	defaultOnPostBeforeReadBodyGetter = getter
}
