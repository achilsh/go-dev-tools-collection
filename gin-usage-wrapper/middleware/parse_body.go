package middleware

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ParseBody() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		requestID := ctx.GetHeader(CtxRequestID)
		if requestID == "" {
			requestID = GenerateRequestID()
		}
		ctx.Writer.Header().Set(CtxRequestID, requestID)
		ctx.Set(CtxRequestID, requestID)
		path := ctx.Request.URL.Path
		method := ctx.Request.Method
		if path == "/test" || path == "/metrics" {
			return
		}

		var reqBuf []byte
		// var reqBuf []byte
		if method == http.MethodPost ||
			method == http.MethodPut ||
			method == http.MethodDelete {
			//

			reqBuf, _ = io.ReadAll(ctx.Request.Body)
			//
		} else if method == http.MethodGet {
			m := map[string]interface{}{}
			for k, v := range ctx.Request.URL.Query() {
				m[k] = v[0]
			}
			reqBuf, _ = json.Marshal(m)
		}

		SetRequestBody(ctx, reqBuf)
		ctx.Next()
	}
}
