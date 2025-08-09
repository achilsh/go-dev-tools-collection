package middleware

import (
	"fmt"
	"net/http"
)

type AuthInterceptorMiddleware struct {
	// 中间件定义
}

func NewAuthInterceptorMiddleware() *AuthInterceptorMiddleware {
	return &AuthInterceptorMiddleware{}
}

func (m *AuthInterceptorMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO generate middleware implement function, delete after code implementation
		fmt.Println("call before.")
		// Passthrough to next handler if need
		next(w, r)
		fmt.Println("call end...")
	}
}
