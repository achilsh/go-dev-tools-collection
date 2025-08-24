package middleware

import (
	"context"
	"fmt"
	"net/http"
)

// 自定义的中间件
/**
 * Func: HeaderMiddlewares is for ...
 *
 * @author GoCommnets
 *
 * @params ...
 * @return
 */

func HeaderMiddlewares(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		xabcHeader := r.Header.Get("x-abc")
		fmt.Println("xabcHeader: ", xabcHeader)

		ctx := context.WithValue(r.Context(), "x-abc", xabcHeader)

		// 创建新的请求对象，使用包含x-abc信息的新context
		r = r.WithContext(ctx)
		// fmt.Println("xabHeader: ", xabcHeader)
		next(w, r)
	}
}
