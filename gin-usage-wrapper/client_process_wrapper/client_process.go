package clientprocesswrapper

import (
	// "content_svr/pub/errors"
	// "content_svr/pub/logger"
	// "content_svr/pub/middleware"
	// "content_svr/pub/requestid"
	// "content_svr/pub/utils"
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"

	errordef "github.com/achilsh/go-dev-tools-collection/gin-usage-wrapper/error_def"
	"github.com/achilsh/go-dev-tools-collection/gin-usage-wrapper/middleware"
	"github.com/gin-gonic/gin"
)

// form accesss client.
func WrapFormClient(handler any) func(ctx *gin.Context) {
	hType := reflect.TypeOf(handler)
	hValue := reflect.ValueOf(handler)
	realFc := func(ctx *gin.Context) {}

	realFuncImpl := func(args []reflect.Value) (results []reflect.Value) {
		ctx := args[0].Interface().(*gin.Context)

		beginTime := time.Now().UnixNano() / 1e6
		nowTm := time.Now().UnixMilli()
		responseData := &errordef.CliError{
			Message:   "",
			MessageId: middleware.GetRequestID(ctx),
			Status:    1000,
			Timestamp: nowTm,
			CostTime:  nowTm - beginTime,
		}

		var realIN []reflect.Value
		realIN = append(realIN, args[0])

		if hType.NumIn() == 2 {
			param := hType.In(1)
			if param.Kind() == reflect.Ptr {
				param = param.Elem()
			}

			if ctx.Request == nil {
				log.Printf("gin receive request is nil.")
				responseData.Status = 4001
				responseData.Message = "request is nil"
				ctx.JSON(http.StatusOK, responseData)
				ctx.Abort()
				return nil
			}

			var fileFieldNames = make(map[string]string)
			mform, err := ctx.MultipartForm()
			if err != nil {
				log.Printf("get multi part form fail, err: %v", err)
				responseData.Status = 4001
				responseData.Message = "request is nil"
				ctx.JSON(http.StatusOK, responseData)
				ctx.Abort()
				return nil
			}

			log.Printf("multipartform value: %+v", mform)
			if mform != nil {
				for fileField, fileFieldValue := range mform.File {
					if fileField != "" {
						if len(fileFieldValue) > 0 {
							if fileFieldValue[0] != nil {
								fileFieldNames[fileField] = fileFieldValue[0].Filename
							}
						}
					}
				}
			}

			log.Printf("fieldNames: %+v", fileFieldNames)
			//
			var fileFieldValue = make(map[string]*bytes.Buffer)
			for fileFieldName := range fileFieldNames {
				file, _, err := ctx.Request.FormFile(fileFieldName)
				if err != nil {
					continue
				}

				var buf *bytes.Buffer = new(bytes.Buffer)
				_, err = io.Copy(buf, file)
				if err != nil {
					file.Close()
					log.Printf("read from file form fail, err: %v")
					continue
				}
				file.Close()
				fileFieldValue[fileFieldName] = buf
			}

			//
			val := reflect.New(param)

			if ctx.Request.Method == http.MethodGet {
				if err := ctx.ShouldBindQuery(val.Interface()); err != nil {
					log.Printf("bind to query failed. err=%v", err)
					ctx.AbortWithStatus(http.StatusInternalServerError)
					return nil
				}
			} else {
				v2 := val.Elem()
				t2 := v2.Type()

				var fileName string = ""
				for i := 0; i < t2.NumField(); i++ {
					field := t2.Field(i)
					tag := field.Tag.Get("form")
					if tag == "" {
						tag = field.Name
					}
					fieldValue := v2.Field(i)
					if !fieldValue.CanSet() {
						continue
					}

					_, ok := fileFieldNames[tag]
					if ok {
						log.Printf("file field name: %v", tag)
						fieldValue.SetBytes(fileFieldValue[tag].Bytes())
						fileName = fileFieldNames[tag]
						continue
					}

					formValue := ctx.Request.FormValue(tag)
					if formValue == "" {
						//内部参数，硬编码
						if tag == "FileName" {
							fieldValue.SetString(fileName)
						}
						continue
					}

					switch field.Type.Kind() {
					case reflect.String:
						log.Printf("parse form key: %v, value: %v", tag, formValue)
						fieldValue.SetString(formValue)

					case reflect.Int:
						intVal, err := strconv.Atoi(formValue)
						if err != nil {
							return
						}
						fieldValue.SetInt(int64(intVal))

					default:
						log.Printf("unhandled default case")
					}
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
			log.Printf("<====InLog: %v, http body: %+v", callFuncName, realIN[1].Interface())
		}

		vals := hValue.Call(realIN)
		statusCode := http.StatusOK

		// 最后一个返回参数是error
		valNum := hType.NumOut()
		if valNum == 2 || valNum == 1 {
			if vals[valNum-1].Interface() != nil {
				log.Printf("resp failed. err=%v", vals[valNum-1].Interface())
				//应答失败
				errImpl, ok := vals[valNum-1].Interface().(*errordef.CliError)
				if ok {
					responseData = errImpl
				} else {
					log.Printf("internal server", vals[valNum-1].Interface().(error))
				}
			} else {
				//应答成功。
				if valNum == 2 {
					responseData.Content = vals[0].Interface()
				}
			}
			if !ctx.IsAborted() {
				ctx.Set("ctx_status", "success")
				ctx.JSON(statusCode, responseData)
			}
		} else {
			if !ctx.IsAborted() {
				ctx.JSON(http.StatusOK, vals[0].Interface())
			}
		}
		ctx.Next()
		return nil
	}
	h := reflect.MakeFunc(reflect.TypeOf(realFc), realFuncImpl)
	return h.Interface().(func(ctx *gin.Context))

}

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
