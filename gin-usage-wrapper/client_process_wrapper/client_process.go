package clientprocesswrapper

import (
	// "content_svr/pub/errors"
	// "content_svr/pub/logger"
	// "content_svr/pub/middleware"
	// "content_svr/pub/requestid"
	// "content_svr/pub/utils"
	"encoding/json"
	"log"
	"net/http"
	"reflect"
	"strings"
	"time"

	errordef "github.com/achilsh/go-dev-tools-collection/gin-usage-wrapper/error_def"
	"github.com/achilsh/go-dev-tools-collection/gin-usage-wrapper/middleware"
	"github.com/gin-gonic/gin"
)

// ClientProcessWrapper 客户端消息统一处理;封装了各类型消息的处理
func ClientProcessWrapper(handler interface{}) func(ctx *gin.Context) {
	hType := reflect.TypeOf(handler)
	hValue := reflect.ValueOf(handler)
	realFc := func(ctx *gin.Context) {}
	realBody := func(args []reflect.Value) (results []reflect.Value) {
		beginTime := time.Now().UnixMilli()

		ctx := args[0].Interface().(*gin.Context)
		var realIN []reflect.Value
		realIN = append(realIN, args[0])
		//
		if hType.NumIn() == 2 {
			param := hType.In(1)
			if param.Kind() == reflect.Ptr {
				param = param.Elem()
			}
			val := reflect.New(param)
			if ctx.Request.Method == http.MethodGet {
				if err := ctx.ShouldBindQuery(val.Interface()); err != nil {
					log.Printf("bind to query failed. err=%v", err)
					ctx.AbortWithStatus(http.StatusInternalServerError)
					return nil
				}
			} else {
				body := middleware.GetRequestBody(ctx)
				if len(body) != 0 {
					if err := json.Unmarshal(body, val.Interface()); err != nil {
						log.Printf("json Unmarshal to struct failed. data=%v, err=%v",
							string(middleware.GetRequestBody(ctx)), err)
						//
						ctx.AbortWithStatus(http.StatusInternalServerError)
						return nil
					}
				}

			}
			for _i := 0; _i < val.Elem().NumField(); _i++ {
				if val.Elem().Field(_i).Kind() == reflect.String {
					val.Elem().Field(_i).SetString(strings.Trim(val.Elem().Field(_i).String(), " "))
				}
			}
			realIN = append(realIN, val)
		}
		vals := hValue.Call(realIN)
		// 返回的格式是: outData, error
		nowTm := time.Now().UnixMilli()

		valNum := hType.NumOut()
		var retMessage *errordef.CliError = &errordef.CliError{
			Message:   "",
			MessageId: middleware.GetRequestID(ctx),
			Status:    1000,
			Timestamp: nowTm,
			CostTime:  nowTm - beginTime,
		}

		statusCode := http.StatusOK
		if valNum == 2 || valNum == 1 {
			if vals[valNum-1].Interface() != nil {
				log.Printf("resp failed. err=%v", vals[valNum-1].Interface())
				//应答失败
				errImpl, ok := vals[valNum-1].Interface().(*errordef.CliError)
				if ok {
					retMessage = errImpl
				} else {
					log.Printf("internal server", vals[valNum-1].Interface().(error))
				}
			} else {
				//应答成功。
				if valNum == 2 {
					retMessage.Content = vals[0].Interface()
				}
			}
			if !ctx.IsAborted() {
				ctx.Set("ctx_status", "success")
				ctx.JSON(statusCode, retMessage)
			}
		} else {
			if !ctx.IsAborted() {
				ctx.JSON(http.StatusOK, vals[0].Interface())
			}
		}
		ctx.Next()
		return nil
	}
	h := reflect.MakeFunc(reflect.TypeOf(realFc), realBody)
	return h.Interface().(func(ctx *gin.Context))
}
