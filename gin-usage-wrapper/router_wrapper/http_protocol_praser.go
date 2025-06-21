package router_wrapper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/bytedance/sonic"
	"github.com/gin-gonic/gin"

	"github.com/achilsh/go-dev-tools-collection/gin-usage-wrapper/error_def"
	"github.com/achilsh/go-dev-tools-collection/gin-usage-wrapper/middleware"

	logger "github.com/achilsh/go-dev-tools-collection/base-lib/log"
)

func MarshalOpFormValue(dstMap map[string]any, src *multipart.Form) {
	// 解析普通表单字段
	for key, values := range src.Value {
		if len(values) == 1 {
			dstMap[key] = values[0]
		} else {
			dstMap[key] = values[0]
		}
	}
	logger.Debugf("data: %+v", dstMap)
}

// form accesss client.
func WrapFormClient(
	handler any,
	beforeRecvHandles ...func(*gin.Context) (int, error),
) func(ctx *gin.Context) {
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
		hasForm := false

		if hType.NumIn() == 2 {
			param := hType.In(1)
			if param.Kind() == reflect.Ptr {
				param = param.Elem()
			}

			if ctx.Request == nil {
				logger.ErrorfCtx(ctx, "gin receive request is nil.")
				responseData.ErrorCode = "4001"
				responseData.ErrorMessage = "request is nil"
				ctx.JSON(http.StatusOK, responseData)
				ctx.Abort()
				return nil
			}
			mform, err := ctx.MultipartForm()
			if err != nil {
				logger.ErrorfCtx(ctx, "get multi part form fail, err: %v", err)
				responseData.ErrorCode = "4001"
				responseData.ErrorMessage = "request is nil"
				ctx.JSON(http.StatusOK, responseData)
				ctx.Abort()
				return nil
			}

			if mform != nil {
				for _, handleFunc := range beforeRecvHandles {
					if handleFunc == nil {
						continue
					}
					if retInt, err := handleFunc(ctx); err != nil {
						logger.ErrorfCtx(ctx, "handle biz logic fail, err: %v", err)
						responseData.ErrorCode = fmt.Sprintf("%v", retInt)
						responseData.ErrorMessage = err.Error()
						ctx.JSON(http.StatusOK, responseData)
						ctx.Abort()
						return nil
					}
				}
			}

			var fileNameContentsMap = make(map[string]map[string]*bytes.Buffer)
			// logger.Debugf("multipartform value: %+v", mform)
			if mform != nil {
				//目前支持多个文件上传只使用 一个 表单字段名
				for fileField, fileFieldValue := range mform.File {
					fileListMap, ok := fileNameContentsMap[fileField]
					if !ok {
						fileListMap = make(map[string]*bytes.Buffer)
						fileNameContentsMap[fileField] = fileListMap
					}

					if len(fileField) > 0 && len(fileFieldValue) > 0 {
						wg := sync.WaitGroup{}
						chData := make(chan any, len(fileFieldValue))

						for ii := range fileFieldValue {
							fv := fileFieldValue[ii]
							if fv != nil {
								wg.Add(1)

								go func(fvIn *multipart.FileHeader) {
									defer wg.Done()

									fH, err := fvIn.Open()
									if err != nil {
										return
									}
									defer fH.Close()

									contentBuf := new(bytes.Buffer)
									_, err = io.Copy(contentBuf, fH)
									if err != nil {
										logger.WarnfCtx(
											ctx,
											"read file fail, fileName: %v, err: %v",
											fvIn.Filename,
											err,
										)
										return
									}

									chData <- struct {
										Name    string
										Content *bytes.Buffer
									}{
										fvIn.Filename,
										contentBuf,
									}

								}(fv)
							}
						}

						wg.Wait()
						close(chData)

						var stopForReceiveOp = false
						for !stopForReceiveOp {
							select {
							case data, isOn := <-chData:
								if !isOn {
									logger.InfofCtx(ctx, "; receive close data channel.")
									stopForReceiveOp = true
									break
								}

								// logger.Debugf("; recevie file content: %+v", data)
								bizData, isType := data.(struct {
									Name    string
									Content *bytes.Buffer
								})
								if isType {
									fileListMap[bizData.Name] = bizData.Content
								}
							}
						}

						if len(fileListMap) > 0 {
							logger.InfofCtx(ctx, "; set file tags: %v", fileField)
							fileNameContentsMap[fileField] = fileListMap
							hasForm = true
						}
					}
				}
			}

			val := reflect.New(param)
			if ctx.Request.Method == http.MethodGet {
				if err := ctx.ShouldBindQuery(val.Interface()); err != nil {
					logger.ErrorfCtx(ctx, "bind to query failed. err=%v", err)
					ctx.AbortWithStatus(http.StatusInternalServerError)
					return nil
				}
			} else {
				v2 := val.Elem()
				t2 := v2.Type()

				// var fileName string = ""
				for i := 0; i < t2.NumField(); i++ {
					// type of field of struct
					field := t2.Field(i)
					tag := field.Tag.Get("form")
					if tag == "" {
						tag = field.Name
					}

					// value of field of struct
					fieldValue := v2.Field(i)
					if !fieldValue.CanSet() {
						continue
					}

					// 其中的一个文件列表；该 tag 对应一个文件列表， field 是 map 表示的文件列表
					fileListContent, ok := fileNameContentsMap[tag]
					if ok {
						logger.DebugfCtx(ctx, "file field name: %v", tag)
						if len(fileListContent) <= 0 {
							continue
						}

						if field.Type.Kind() == reflect.Map {
							if fieldValue.IsNil() { //如果是个空的 map
								fieldValue.Set(reflect.MakeMap(fieldValue.Type()))
							}

							for formItemK, formItemV := range fileListContent {
								fieldValue.SetMapIndex(reflect.ValueOf(formItemK), reflect.ValueOf(formItemV.Bytes()))
							}
						}
						continue
					}

					formValue := ctx.Request.FormValue(tag)
					if formValue == "" {
						continue
					}

					switch field.Type.Kind() {
					case reflect.String:
						logger.InfofCtx(ctx, "parse form key: %v, value: %v", tag, formValue)
						fieldValue.SetString(formValue)

					case reflect.Int, reflect.Uint,
						reflect.Uint16, reflect.Int16,
						reflect.Int32, reflect.Uint32,
						reflect.Int64, reflect.Uint64,
						reflect.Int8, reflect.Uint8:

						intVal, err := strconv.ParseInt(formValue, 10, 64)
						if err != nil {
							logger.ErrorfCtx(ctx, "parse int form field, err: %v", err)
							continue
						}
						fieldValue.SetInt(int64(intVal))

					case reflect.Float32, reflect.Float64:
						f64, err := strconv.ParseFloat(formValue, 64)
						if err != nil {
							logger.ErrorfCtx(ctx, "parse float form field, err: %v", err)
							continue
						}
						fieldValue.SetFloat(f64)

					default:
						logger.ErrorfCtx(ctx, "unhandled default case")
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

			if hasForm {
				logger.AccessfCtx(ctx, "=======> InPutLog: %v", callFuncName)
			} else {
				logger.AccessfCtx(ctx, "=======> InPutLog: %v, http body: %+v", callFuncName, realIN[1].Interface())
			}
		}

		vals := hValue.Call(realIN)

		// 最后一个返回参数是error
		valOutNum := hType.NumOut()

		if valOutNum == 2 { //返回两个字段
			statusCode := http.StatusOK

			if vals[valOutNum-1].Interface() != nil { //返回非nil的错误
				logger.ErrorfCtx(
					ctx,
					"resp failed. err=%v",
					error_def.StructToJsonString(vals[valOutNum-1].Interface()),
				)
				//应答失败
				errImpl, ok := vals[valOutNum-1].Interface().(error_def.CliErrorEr)
				if ok {
					responseData.ErrorCode = errImpl.GetCode()
					responseData.ErrorMessage = errImpl.GetCodeMsg()
				} else {
					logger.ErrorfCtx(ctx, "no known error: %+v", vals[valOutNum-1].Interface().(error))
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

				respJson, _ := sonic.Marshal(responseData)
				if !IsNoOutputLog(ctx) {
					logger.AccessfCtx(ctx, "<-------- %v, http response: %+v", callFuncName, string(respJson))
				}
				ctx.JSON(statusCode, responseData)
			}
		} else if valOutNum == 1 { // 只返回一个参数，没有返回具体业务的数据
			if vals[valOutNum-1].Interface() != nil { //返回非nil的错误
				logger.ErrorfCtx(ctx, "resp failed. err=%v", error_def.StructToJsonString(vals[valOutNum-1].Interface()))
				//应答失败
				errImpl, ok := vals[valOutNum-1].Interface().(error_def.CliErrorEr)
				if ok {
					responseData.ErrorCode = errImpl.GetCode()
					responseData.ErrorMessage = errImpl.GetCodeMsg()
				} else {
					logger.ErrorfCtx(ctx, "no known error: %+v", vals[valOutNum-1].Interface().(error))
					responseData.ErrorCode = "5000"
					responseData.ErrorMessage = "unknown error in server."
				}
			} else {
				//应答成功。
				responseData.ErrorCode = "200"
				responseData.ErrorMessage = "success"
				responseData.CostTimeMs = time.Now().UnixNano()/1e6 - beginTime
			}

			if !ctx.IsAborted() {
				logger.AccessfCtx(ctx, "OutLog: %v, http response: %+v", callFuncName, responseData)
				ctx.JSON(http.StatusOK, responseData)
			}
		} else {
			logger.ErrorfCtx(ctx, "return out parameter nums is more than 2.")
		}
		ctx.Next()
		return nil
	}
	h := reflect.MakeFunc(reflect.TypeOf(realFc), realFuncImpl)
	return h.Interface().(func(ctx *gin.Context))

}

// WrapperClient 包装 gin 解包和处理
func WrapperClient(
	handler interface{},
	beforeReadFunc ...func(ctx *gin.Context) (int, error),
) func(ctx *gin.Context) {
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
					logger.ErrorfCtx(ctx, "bind to query failed. err=%v", err)
					ctx.AbortWithStatus(http.StatusInternalServerError)
					return nil
				}
			} else {
				body := middleware.GetRequestBody(ctx)
				if len(body) > 0 {
					for _, bHandleFunc := range beforeReadFunc {
						if bHandleFunc == nil {
							continue
						}
						if retInt, err := bHandleFunc(ctx); err != nil {
							logger.ErrorfCtx(ctx, "before handler fail, err: %v", err)
							responseData.ErrorCode = fmt.Sprintf("%v", retInt)
							responseData.ErrorMessage = err.Error()
							ctx.JSON(http.StatusOK, responseData)
							ctx.Abort()
							return nil
						}
					}
				}
				if len(body) != 0 {
					if err := json.Unmarshal(body, val.Interface()); err != nil {
						logger.ErrorfCtx(ctx, "json Unmarshal to struct failed. data=%v, err=%v", body, err)
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
			//
			realIN = append(realIN, val)
		}

		callFuncName := ""
		if len(realIN) >= 2 {
			lastFuncNames := strings.Split(runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name(), ".")
			if len(lastFuncNames) > 0 {
				callFuncName = lastFuncNames[len(lastFuncNames)-1]
			}
			logger.AccessfCtx(ctx, "===>InLog: %v, http body: %+v", callFuncName, realIN[1].Interface())
		}

		vals := hValue.Call(realIN)

		// 最后一个返回参数是error
		valOutNum := hType.NumOut()

		if valOutNum == 2 { //返回两个字段
			statusCode := http.StatusOK

			if vals[valOutNum-1].Interface() != nil { //返回非nil的错误
				logger.ErrorfCtx(
					ctx,
					"resp failed. err=%v",
					error_def.StructToJsonString(vals[valOutNum-1].Interface()),
				)
				//应答失败
				errImpl, ok := vals[valOutNum-1].Interface().(error_def.CliErrorEr)
				if ok {
					responseData.ErrorCode = errImpl.GetCode()
					responseData.ErrorMessage = errImpl.GetCodeMsg()
				} else {
					logger.ErrorfCtx(ctx, "no known error: %+v", vals[valOutNum-1].Interface().(error))
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

				respJson, _ := sonic.Marshal(responseData)
				if !IsNoOutputLog(ctx) {
					logger.AccessfCtx(ctx, "<-------- %v, http response: %+v", callFuncName, string(respJson))
				}
				ctx.JSON(statusCode, responseData)
			}
		} else if valOutNum == 1 { // 只返回一个参数，没有返回具体业务的数据
			if vals[valOutNum-1].Interface() != nil { //返回非nil的错误
				logger.ErrorfCtx(ctx, "resp failed. err=%v", error_def.StructToJsonString(vals[valOutNum-1].Interface()))
				//应答失败
				errImpl, ok := vals[valOutNum-1].Interface().(error_def.CliErrorEr)
				if ok {
					responseData.ErrorCode = errImpl.GetCode()
					responseData.ErrorMessage = errImpl.GetCodeMsg()
				} else {
					logger.ErrorfCtx(ctx, "no known error: %+v", vals[valOutNum-1].Interface().(error))
					responseData.ErrorCode = "5000"
					responseData.ErrorMessage = "unknown error in server."
				}
			} else {
				//应答成功。
				responseData.ErrorCode = "200"
				responseData.ErrorMessage = "success"
				responseData.CostTimeMs = time.Now().UnixNano()/1e6 - beginTime
			}

			if !ctx.IsAborted() {
				respJson, _ := sonic.Marshal(responseData)
				if !IsNoOutputLog(ctx) {
					logger.AccessfCtx(ctx, "<-------- %v, http response: %+v", callFuncName, string(respJson))
				}
				ctx.JSON(http.StatusOK, responseData)
			}
		} else {
			logger.ErrorfCtx(ctx, "return out parameter nums is more than 2.")
		}
		ctx.Next()
		return nil
	}
	h := reflect.MakeFunc(reflect.TypeOf(realFc), realFuncImpl)
	return h.Interface().(func(ctx *gin.Context))
}
