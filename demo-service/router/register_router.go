package router

import (
	"github.com/achilsh/go-dev-tools-collection/demo-service/handler"
	rw "github.com/achilsh/go-dev-tools-collection/gin-usage-wrapper/router_wrapper"

	"github.com/achilsh/go-dev-tools-collection/demo-service/service/middleware"
	"github.com/gin-gonic/gin"
)

const (
	UrlPathSuffixWelcomeWord = "demo-call"
)

// 设置读取消息体（未解码）之后， 调用处理业务逻辑函数之前的默认处理逻辑。： 比如body 体鉴权等。
func RegisterBeforeReadHandle() {
	rw.SetBeforeReadBodyOnFormHandlerGetter(
		func() []func(ctx *gin.Context) (int, error) {
			return []func(ctx *gin.Context) (int, error){middleware.CheckBeforeReadBodyOnForm}
		})

	rw.SetBeforeReadBodyPostHandlerGetter(
		func() []func(ctx *gin.Context) (int, error) {
			return []func(ctx *gin.Context) (int, error){middleware.CheckBeforeReadBodyOnJson}
		})
}
func RegisterRouter() {

	RegisterBeforeReadHandle()

	srvRouteV1 := GetRouter().Group("/demo-server/v1")
	{
		//底层使用反射来接收数据
		rw.RegisterPostNoInProcess(srvRouteV1, "x1", handler.DemoIn)
		//底层使用泛型来接收数据
		rw.RegisterPostGenericNoInProcess(srvRouteV1, "x1/generic", handler.DemoInPtr)

		//底层使用反射来接收数据
		rw.RegisterPostProcess(srvRouteV1, "x2", handler.DemoInOut)
		//底层使用泛型来接收数据
		rw.RegisterPostGenericProcess(srvRouteV1, "x2/generic", handler.DemoABC)
		//	底层使用反射来接收数据
		rw.RegisterPostProcess(srvRouteV1, "z1/abc", handler.DemoABC)

		//底层使用反射来接收数据
		rw.RegisterPostNoOutProcess(srvRouteV1, "x3", handler.DemoNoOut)
		//底层使用反射来接收数据
		rw.RegisterPostGenericNoOutProcess(srvRouteV1, "x3/generic", handler.DemoNoOut)

		//底层使用反射来接收数据
		rw.RegisterGetProcess(srvRouteV1, "y1", handler.DemoGetInOut)
		//底层使用反射来接收数据
		rw.RegisterGetGenericProcess(srvRouteV1, "y1/generic", handler.DemoGetInOut)
		//底层使用反射来接收数据
		rw.RegisterGetNoInProcess(srvRouteV1, "y2", handler.DemoGetNoIn)
		//底层使用反射来接收数据
		rw.RegisterGetNoInGenericProcess(srvRouteV1, "y2/generic", handler.DemoGetNoInPtr)
		//底层使用反射来接收数据
		rw.RegisterGetNoOutProcess(srvRouteV1, "y3", handler.DemoGetNoOut)
		//底层使用反射来接收数据
		rw.RegisterGetNoOutGenericProcess(srvRouteV1, "y3/generic", handler.DemoGetNoOut)
	}
	{
		//底层使用反射来接收数据
		rw.RegisterPostForm(srvRouteV1, "f1", handler.DemoFormInOut)
		//底层使用反射来接收数据
		rw.RegisterPostGenericForm(srvRouteV1, "f1/generic", handler.DemoFormInOut)
	}

	//底层使用反射来接收数据
	rw.RegisterPostProcess(srvRouteV1, UrlPathSuffixWelcomeWord, handler.WelcomeWordHandle)

}
