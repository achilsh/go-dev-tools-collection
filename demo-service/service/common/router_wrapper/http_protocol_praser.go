package router_wrapper

import (
	"encoding/json"
	"net/http"
	"reflect"
	"runtime"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	//"github.com/mgechev/revive/config"

	"demo-service/service/middleware"
	"demo-service/service/utils/error_def"

	logger "github.com/achilsh/go-dev-tools-collection/base-lib/log"
)

// WrapperClient 包装 gin 解包和处理
func WrapperClient(handler interface{}) func(ctx *gin.Context) {
	hType := reflect.TypeOf(handler)
	hValue := reflect.ValueOf(handler)
	realFc := func(ctx *gin.Context) {}

	realFuncImpl := func(args []reflect.Value) (results []reflect.Value) {
		beginTime := time.Now().UnixNano() / 1e6

		responseData := &error_def.HttpResponse{
			ErrorMessage: "success",
			ErrorCode:    "",
		}

		ctx := args[0].Interface().(*gin.Context)
		var realIN []reflect.Value
		realIN = append(realIN, args[0])

		if hType.NumIn() == 2 {
			param := hType.In(1)
			if param.Kind() == reflect.Ptr {
				param = param.Elem()
			}
			//
			val := reflect.New(param)
			if ctx.Request.Method == http.MethodGet {
				if err := ctx.ShouldBindQuery(val.Interface()); err != nil {
					logger.Errorf("bind to query failed. err=%v", err)
					ctx.AbortWithStatus(http.StatusInternalServerError)
					return nil
				}
			} else {
				body := middleware.GetRequestBody(ctx)
				if len(body) != 0 {
					if err := json.Unmarshal(body, val.Interface()); err != nil {
						logger.Errorf("json Unmarshal to struct failed. data=%v, err=%v", body, err)

						ctx.AbortWithStatus(http.StatusInternalServerError)
						return nil
					}

					var dataBodyMap map[string]any
					json.Unmarshal(body, dataBodyMap)
					errCode, payLoadPtr := middleware.ParseToken(dataBodyMap, ctx)
					if errCode != "" {
						logger.Errorf("accessToken Check interface fail")
						responseData.ErrorCode = errCode
						responseData.ErrorMessage = "accessToken verification failed."
						ctx.JSON(http.StatusOK, responseData)
						ctx.Abort()
						return nil
					}

					middleware.SetPayLoadToCtx(payLoadPtr, ctx)
				}
			}
			for _i := 0; _i < val.Elem().NumField(); _i++ {
				if val.Elem().Field(_i).Kind() == reflect.String {
					val.Elem().Field(_i).SetString(strings.Trim(val.Elem().Field(_i).String(), " "))
				}
			}
			//
			realIN = append(realIN, val)
		}

		callFuncName := ""
		if len(realIN) >= 2 {
			lastFuncNames := strings.Split(runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name(), ".")
			if len(lastFuncNames) > 0 {
				callFuncName = lastFuncNames[len(lastFuncNames)-1]
			}
			logger.Infof("<====InLog: %v, http body: %+v", callFuncName, realIN[1].Interface())
		}

		vals := hValue.Call(realIN)

		// 最后一个返回参数是error
		valOutNum := hType.NumOut()

		if valOutNum == 2 { //返回两个字段
			statusCode := http.StatusOK

			if vals[valOutNum-1].Interface() != nil { //返回非nil的错误
				logger.Errorf("resp failed. err=%v", error_def.StructToJsonString(vals[valOutNum-1].Interface()))
				//应答失败
				errImpl, ok := vals[valOutNum-1].Interface().(error_def.CliErrorEr)
				if ok {
					responseData.ErrorCode = errImpl.GetCode()
					responseData.ErrorMessage = errImpl.GetCodeMsg()
				} else {
					logger.Errorf("no known error: %+v", vals[valOutNum-1].Interface().(error))
					responseData.ErrorCode = "5000"
					responseData.ErrorMessage = "unknown error in server."
				}
			} else {
				//应答成功。
				responseData.ErrorCode = "200"
				responseData.ErrorMessage = "success"
				responseData.CostTimeMs = time.Now().UnixNano()/1e6 - beginTime
				responseData.Data = vals[0].Interface()
			}
			if !ctx.IsAborted() {
				ctx.Set("ctx_status", "success")

				logger.Infof("<-------- %v, http response: %+v", callFuncName, responseData)
				ctx.JSON(statusCode, responseData)
			}
		} else if valOutNum == 1 { // 只返回一个参数，没有返回具体业务的数据
			if vals[valOutNum-1].Interface() != nil { //返回非nil的错误
				logger.Errorf("resp failed. err=%v", error_def.StructToJsonString(vals[valOutNum-1].Interface()))
				//应答失败
				errImpl, ok := vals[valOutNum-1].Interface().(error_def.CliErrorEr)
				if ok {
					responseData.ErrorCode = errImpl.GetCode()
					responseData.ErrorMessage = errImpl.GetCodeMsg()
				} else {
					logger.Errorf("no known error: %+v", vals[valOutNum-1].Interface().(error))
					responseData.ErrorCode = "5000"
					responseData.ErrorMessage = "unknown error in server."
				}
			} else {
				//应答成功。
				responseData.ErrorCode = "200"
				responseData.ErrorMessage = "success"
				responseData.CostTimeMs = time.Now().UnixNano()/1e6 - beginTime
				responseData.Data = vals[0].Interface()
			}

			if !ctx.IsAborted() {
				logger.Infof("OutLog: %v, http response: %+v", callFuncName, responseData)
				ctx.JSON(http.StatusOK, responseData)
			}
		} else {
			logger.Errorf("return out parameter nums is more than 2.")
		}
		ctx.Next()
		return nil
	}
	h := reflect.MakeFunc(reflect.TypeOf(realFc), realFuncImpl)
	return h.Interface().(func(ctx *gin.Context))
}
