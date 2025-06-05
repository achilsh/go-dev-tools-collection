package router_wrapper

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"runtime"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	logger "github.com/achilsh/go-dev-tools-collection/base-lib/log"
	"github.com/achilsh/go-dev-tools-collection/demo-service/service/middleware"
	"github.com/achilsh/go-dev-tools-collection/demo-service/service/utils/error_def"
)

// 既有入参，也有返回值参数
func WrapperGenericClient[C context.Context, IN, OUT any](
	fn func(C, *IN) (*OUT, error),
	beforeReadFunc ...func(*gin.Context) (int, error),
) func(ctx *gin.Context) {
	return func(c *gin.Context) {
		beginTime := time.Now().UnixNano() / 1e6
		responseData := &error_def.HttpResponse{
			ErrorMessage: "success",
			ErrorCode:    "",
		}

		var in IN
		if c.Request.Method == http.MethodGet {
			if err := c.ShouldBindQuery(&in); err != nil {
				logger.ErrorfCtx(c, "bind to query failed. err=%v", err)
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
		} else {
			body := middleware.GetRequestBody(c)
			if len(body) > 0 {
				for _, bHandleFunc := range beforeReadFunc {
					if bHandleFunc == nil {
						continue
					}
					if retInt, err := bHandleFunc(c); err != nil {
						logger.ErrorfCtx(c, "before handler fail, err: %v", err)
						responseData.ErrorCode = fmt.Sprintf("%v", retInt)
						responseData.ErrorMessage = err.Error()
						c.JSON(http.StatusOK, responseData)
						c.Abort()
						return
					}
				}
			}
			if len(body) != 0 {
				if err := json.Unmarshal(body, &in); err != nil {
					logger.ErrorfCtx(c, "json Unmarshal to struct failed. data=%v, err=%v", body, err)
					c.AbortWithStatus(http.StatusInternalServerError)
					return
				}
			}
		}

		callFuncName := ""
		lastFuncNames := strings.Split(runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name(), ".")
		if len(lastFuncNames) > 0 {
			callFuncName = lastFuncNames[len(lastFuncNames)-1]
		}
		logger.AccessfCtx(c, "<====InLog: %v, http body: %+v", callFuncName, in)

		var ctx context.Context = c
		out, err := fn(ctx.(C), &in)
		statusCode := http.StatusOK
		if err != nil { //返回非nil的错误
			logger.ErrorfCtx(
				ctx,
				"resp failed. err=%v",
				error_def.StructToJsonString(err),
			)
			//应答失败
			errImpl, ok := err.(error_def.CliErrorEr)
			if ok {
				responseData.ErrorCode = errImpl.GetCode()
				responseData.ErrorMessage = errImpl.GetCodeMsg()
			} else {
				logger.ErrorfCtx(ctx, "no known error: %+v", err)
				responseData.ErrorCode = "5000"
				responseData.ErrorMessage = "unknown error in server: " + err.Error()
			}
		} else {
			//应答成功。
			responseData.ErrorCode = "200"
			responseData.ErrorMessage = "success"
			responseData.CostTimeMs = time.Now().UnixNano()/1e6 - beginTime
			responseData.Data = out
		}
		if !c.IsAborted() {
			c.Set("ctx_status", "success")

			logger.AccessfCtx(ctx, "<-------- %v, http response: %+v", callFuncName, responseData)
			c.JSON(statusCode, responseData)
		}
		c.Next()
	}
}

func WrapperGenericClientNoIN[C context.Context, OUT any](
	fn func(C) (*OUT, error),
	beforeReadFunc ...func(*gin.Context) (int, error),
) func(ctx *gin.Context) {
	return func(c *gin.Context) {
		beginTime := time.Now().UnixNano() / 1e6
		responseData := &error_def.HttpResponse{
			ErrorMessage: "success",
			ErrorCode:    "",
		}

		callFuncName := ""

		var ctx context.Context = c
		out, err := fn(ctx.(C))
		statusCode := http.StatusOK
		if err != nil { //返回非nil的错误
			logger.ErrorfCtx(
				ctx,
				"resp failed. err=%v",
				error_def.StructToJsonString(err),
			)
			//应答失败
			errImpl, ok := err.(error_def.CliErrorEr)
			if ok {
				responseData.ErrorCode = errImpl.GetCode()
				responseData.ErrorMessage = errImpl.GetCodeMsg()
			} else {
				logger.ErrorfCtx(ctx, "no known error: %+v", err)
				responseData.ErrorCode = "5000"
				responseData.ErrorMessage = "unknown error in server: " + err.Error()
			}
		} else {
			//应答成功。
			responseData.ErrorCode = "200"
			responseData.ErrorMessage = "success"
			responseData.CostTimeMs = time.Now().UnixNano()/1e6 - beginTime
			responseData.Data = out
		}
		if !c.IsAborted() {
			c.Set("ctx_status", "success")

			logger.AccessfCtx(ctx, "<-------- %v, http response: %+v", callFuncName, responseData)
			c.JSON(statusCode, responseData)
		}
		c.Next()
	}
}

func WrapperGenericClientNoOUT[C context.Context, IN any](
	fn func(C, *IN) error,
	beforeReadFunc ...func(*gin.Context) (int, error),
) func(ctx *gin.Context) {

	return func(c *gin.Context) {
		beginTime := time.Now().UnixNano() / 1e6
		responseData := &error_def.HttpResponse{
			ErrorMessage: "success",
			ErrorCode:    "",
		}

		var in IN
		if c.Request.Method == http.MethodGet {
			if err := c.ShouldBindQuery(&in); err != nil {
				logger.ErrorfCtx(c, "bind to query failed. err=%v", err)
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
		} else {
			body := middleware.GetRequestBody(c)
			if len(body) > 0 {
				for _, bHandleFunc := range beforeReadFunc {
					if bHandleFunc == nil {
						continue
					}
					if retInt, err := bHandleFunc(c); err != nil {
						logger.ErrorfCtx(c, "before handler fail, err: %v", err)
						responseData.ErrorCode = fmt.Sprintf("%v", retInt)
						responseData.ErrorMessage = err.Error()
						c.JSON(http.StatusOK, responseData)
						c.Abort()
						return
					}
				}
			}
			if len(body) != 0 {
				if err := json.Unmarshal(body, &in); err != nil {
					logger.ErrorfCtx(c, "json Unmarshal to struct failed. data=%v, err=%v", body, err)
					c.AbortWithStatus(http.StatusInternalServerError)
					return
				}
			}
		}

		callFuncName := ""
		lastFuncNames := strings.Split(runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name(), ".")
		if len(lastFuncNames) > 0 {
			callFuncName = lastFuncNames[len(lastFuncNames)-1]
		}
		logger.AccessfCtx(c, "<====InLog: %v, http body: %+v", callFuncName, in)

		var ctx context.Context = c
		err := fn(ctx.(C), &in)
		statusCode := http.StatusOK
		if err != nil { //返回非nil的错误
			logger.ErrorfCtx(
				ctx,
				"resp failed. err=%v",
				error_def.StructToJsonString(err),
			)
			//应答失败
			errImpl, ok := err.(error_def.CliErrorEr)
			if ok {
				responseData.ErrorCode = errImpl.GetCode()
				responseData.ErrorMessage = errImpl.GetCodeMsg()
			} else {
				logger.ErrorfCtx(ctx, "no known error: %+v", err)
				responseData.ErrorCode = "5000"
				responseData.ErrorMessage = "unknown error in server: " + err.Error()
			}
		} else {
			//应答成功。
			responseData.ErrorCode = "200"
			responseData.ErrorMessage = "success"
			responseData.CostTimeMs = time.Now().UnixNano()/1e6 - beginTime
		}
		if !c.IsAborted() {
			c.Set("ctx_status", "success")

			logger.AccessfCtx(ctx, "<-------- %v, http response: %+v", callFuncName, responseData)
			c.JSON(statusCode, responseData)
		}
		c.Next()
	}

}
