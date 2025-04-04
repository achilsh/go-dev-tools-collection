package router

import (
	"demo-service/handler"
	rw "demo-service/service/common/router_wrapper"
)

const (
	UrlPathSuffixWelcomeWord = "demo-call"
)

func RegisterRouter() {
	srvRouteV1 := GetRouter().Group("/demo-server/v1")
	{
		rw.RegisterPostNoInProcess(srvRouteV1, "x1", handler.DemoIn)

		rw.RegisterPostProcess(srvRouteV1, "x2", handler.DemoInOut)

		rw.RegisterPostNoOutProcess(srvRouteV1, "x3", handler.DemoNoOut)
		//
		rw.RegisterGetProcess(srvRouteV1, "y1", handler.DemoGetInOut)
		rw.RegisterGetNoInProcess(srvRouteV1, "y2", handler.DemoGetNoIn)
		rw.RegisterGetNoOutProcess(srvRouteV1, "y3", handler.DemoGetNoOut)
	}
	{
		rw.RegisterPostForm(srvRouteV1, "f1", handler.DemoFormInOut)
	}

	//
	rw.RegisterPostProcess(srvRouteV1, UrlPathSuffixWelcomeWord, handler.WelcomeWordHandle)

}
