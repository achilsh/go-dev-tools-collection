package svc

import (
	"add_server_on_service_http/internal/config"
	"add_server_on_service_http/internal/middleware"

	"github.com/zeromicro/go-zero/rest"
)

type ServiceContext struct {
	Config          config.Config
	AuthInterceptor rest.Middleware
}

// 可以在 servicecontext.go 里面传递依赖给 logic，比如 mysql, redis 传递给logic 层等
func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:          c,
		AuthInterceptor: middleware.NewAuthInterceptorMiddleware().Handle,
	}
}
