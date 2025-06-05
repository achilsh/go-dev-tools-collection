package router

import (
	"github.com/achilsh/go-dev-tools-collection/demo-service/handler"
	rw "github.com/achilsh/go-dev-tools-collection/demo-service/service/common/router_wrapper"
)

const (
	UrlPathSuffixWelcomeWord = "demo-call"
)

func RegisterRouter() {
	srvRouteV1 := GetRouter().Group("/demo-server/v1")
	{
		rw.RegisterPostNoInProcess(srvRouteV1, "x1", handler.DemoIn)
		rw.RegisterPostGenericNoInProcess(srvRouteV1, "x1/generic", handler.DemoInPtr)

		rw.RegisterPostProcess(srvRouteV1, "x2", handler.DemoInOut)
		rw.RegisterPostGenericProcess(srvRouteV1, "x2/generic", handler.DemoABC)
		//
		rw.RegisterPostProcess(srvRouteV1, "z1/abc", handler.DemoABC)

		rw.RegisterPostNoOutProcess(srvRouteV1, "x3", handler.DemoNoOut)
		rw.RegisterPostGenericNoOutProcess(srvRouteV1, "x3/generic", handler.DemoNoOut)
		//
		rw.RegisterGetProcess(srvRouteV1, "y1", handler.DemoGetInOut)
		rw.RegisterGetGenericProcess(srvRouteV1, "y1/generic", handler.DemoGetInOut)
		//
		rw.RegisterGetNoInProcess(srvRouteV1, "y2", handler.DemoGetNoIn)
		rw.RegisterGetNoInGenericProcess(srvRouteV1, "y2/generic", handler.DemoGetNoInPtr)
		//
		rw.RegisterGetNoOutProcess(srvRouteV1, "y3", handler.DemoGetNoOut)
		rw.RegisterGetNoOutGenericProcess(srvRouteV1, "y3/generic", handler.DemoGetNoOut)
	}
	{
		rw.RegisterPostForm(srvRouteV1, "f1", handler.DemoFormInOut)
		//
		rw.RegisterPostGenericForm(srvRouteV1, "f1/generic", handler.DemoFormInOut)
	}

	//
	rw.RegisterPostProcess(srvRouteV1, UrlPathSuffixWelcomeWord, handler.WelcomeWordHandle)

}
