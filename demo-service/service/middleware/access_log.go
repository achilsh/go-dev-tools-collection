package middleware

import (
	"encoding/json"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/achilsh/go-dev-tools-collection/gin-usage-wrapper/middleware"
	"github.com/gin-gonic/gin"
)

const (
	CtxClientIp = "ctx_client_ip"
)

func AccessLogMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set(CtxClientIp, ctx.Request.Header.Get("X-REAL-IP"))
		requestID := ctx.GetHeader(CtxRequestID)
		if requestID == "" {
			requestID = GenerateRequestID()
		}

		ctx.Writer.Header().Set(CtxRequestID, requestID)
		ctx.Set(CtxRequestID, requestID)

		method := ctx.Request.Method
		contentType := ctx.Request.Header.Get("Content-Type")
		contentType = strings.ToLower(contentType)
		regHandle, err := regexp.Compile("multipart|form-data")
		if err == nil && regHandle != nil {
			if regHandle.MatchString(contentType) {
				ctx.Next() //执行后面的中间件和handler，完了再实行下面的。
				return
			}
		}

		var reqBuf []byte
		if method == http.MethodPost || method == http.MethodPut || method == http.MethodDelete {
			reqBuf, _ = io.ReadAll(ctx.Request.Body)

		} else if method == http.MethodGet {
			m := map[string]interface{}{}
			for k, v := range ctx.Request.URL.Query() {
				m[k] = v[0]
			}
			reqBuf, _ = json.Marshal(m)
		}
		middleware.SetRequestBody(ctx, reqBuf)

		ctx.Next() //执行后面的中间件和handler，完了再实行下面的。
	}
}
