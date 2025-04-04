package router_wrapper

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	//"github.com/mgechev/revive/config"

	"demo-service/service/middleware"
	"demo-service/service/utils/error_def"

	logger "github.com/achilsh/go-dev-tools-collection/base-lib/log"
)

// form accesss client.
func WrapFormClient(handler any) func(ctx *gin.Context) {
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

			if ctx.Request == nil {
				logger.Errorf("gin receive request is nil.")
				responseData.ErrorCode = "4001"
				responseData.ErrorMessage = "request is nil"
				ctx.JSON(http.StatusOK, responseData)
				ctx.Abort()
				return nil
			}

			var fileFieldNames = make(map[string]string)
			// fh, err := ctx.FormFile("audio")
			// if err != nil {
			// 	logger.Errorf("form file audio fail, err: %v", err)
			// } else {
			// 	logger.Infof("file header: %v", fh.Filename)
			// }
			mform, err := ctx.MultipartForm()
			if err != nil {
				logger.Errorf("get multi part form fail, err: %v", err)
				responseData.ErrorCode = "4001"
				responseData.ErrorMessage = "request is nil"
				ctx.JSON(http.StatusOK, responseData)
				ctx.Abort()
				return nil
			}

			logger.Debugf("multipartform value: %+v", mform)
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

			logger.Debugf("fieldNames: %+v", fileFieldNames)
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
					logger.Errorf("read from file form fail, err: %v")
					continue
				}
				file.Close()
				fileFieldValue[fileFieldName] = buf
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
						logger.Infof("file field name: %v", tag)
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
						logger.Infof("parse form key: %v, value: %v", tag, formValue)
						fieldValue.SetString(formValue)

					case reflect.Int:
						intVal, err := strconv.Atoi(formValue)
						if err != nil {
							return
						}
						fieldValue.SetInt(int64(intVal))

					default:
						logger.Errorf("unhandled default case")
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
