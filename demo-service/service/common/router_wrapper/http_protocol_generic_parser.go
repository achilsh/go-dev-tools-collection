package router_wrapper

import (
	"bytes"
	"context"
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

			respJson, _ := sonic.Marshal(responseData)
			if !IsNoOutputLog(c) {
				logger.AccessfCtx(ctx, "<-------- %v, http response: %+v", callFuncName, string(respJson))
			}

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

			respJson, _ := sonic.Marshal(responseData)
			if !IsNoOutputLog(c) {
				logger.AccessfCtx(ctx, "<-------- %v, http response: %+v", callFuncName, string(respJson))
			}
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

			respJson, _ := sonic.Marshal(responseData)
			if !IsNoOutputLog(c) {
				logger.AccessfCtx(ctx, "<-------- %v, http response: %+v", callFuncName, string(respJson))
			}
			c.JSON(statusCode, responseData)
		}
		c.Next()
	}

}

// 封装 form with generic interface
func WrapperGenericFormClient[C context.Context, IN, OUT any](
	fn func(C, *IN) (*OUT, error),
	beforeReceiveHandles ...func(*gin.Context) (int, error),
) func(ctx *gin.Context) {

	return func(c *gin.Context) {
		beginTime := time.Now().UnixNano() / 1e6
		responseData := &error_def.HttpResponse{
			ErrorMessage: "success",
			ErrorCode:    "",
		}

		var in IN
		hasForm := false

		if c.Request == nil {
			logger.ErrorfCtx(c, "gin receive request is nil.")
			responseData.ErrorCode = "4001"
			responseData.ErrorMessage = "request is nil"
			c.JSON(http.StatusOK, responseData)
			c.Abort()
			return
		}

		mform, err := c.MultipartForm()
		if err != nil {
			logger.ErrorfCtx(c, "get multi part form fail, err: %v", err)
			responseData.ErrorCode = "4001"
			responseData.ErrorMessage = "request is nil"
			c.JSON(http.StatusOK, responseData)
			c.Abort()
			return
		}

		if mform != nil {
			for _, handleFunc := range beforeReceiveHandles {
				if handleFunc == nil {
					continue
				}
				if retInt, err := handleFunc(c); err != nil {
					logger.ErrorfCtx(c, "handle biz logic fail, err: %v", err)
					responseData.ErrorCode = fmt.Sprintf("%v", retInt)
					responseData.ErrorMessage = err.Error()
					c.JSON(http.StatusOK, responseData)
					c.Abort()
					return
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
										c,
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
								logger.InfofCtx(c, "; receive close data channel.")
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
						logger.InfofCtx(c, "; set file tags: %v", fileField)
						fileNameContentsMap[fileField] = fileListMap
						hasForm = true
					}
				}
			}
		}

		if c.Request.Method == http.MethodGet {
			if err := c.ShouldBindQuery(&in); err != nil {
				logger.ErrorfCtx(c, "bind to query failed. err=%v", err)
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
		} else {
			typeIn := reflect.TypeOf(in)
			valueIn := reflect.ValueOf(&in).Elem()

			for i := 0; i < typeIn.NumField(); i++ {
				// type of field of struct
				field := typeIn.Field(i)
				tag := field.Tag.Get("form")
				if tag == "" {
					tag = field.Name
				}

				// value of field of struct
				fieldValue := valueIn.Field(i)
				if !fieldValue.CanSet() {
					continue
				}

				// 其中的一个文件列表；该 tag 对应一个文件列表， field 是 map 表示的文件列表
				fileListContent, ok := fileNameContentsMap[tag]
				if ok {
					logger.DebugfCtx(c, "file field name: %v", tag)
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

				formValue := c.Request.FormValue(tag)
				if formValue == "" {
					continue
				}

				switch field.Type.Kind() {
				case reflect.String:
					logger.InfofCtx(c, "parse form key: %v, value: %v", tag, formValue)
					fieldValue.SetString(formValue)

				case reflect.Int, reflect.Uint,
					reflect.Uint16, reflect.Int16,
					reflect.Int32, reflect.Uint32,
					reflect.Int64, reflect.Uint64,
					reflect.Int8, reflect.Uint8:

					intVal, err := strconv.ParseInt(formValue, 10, 64)
					if err != nil {
						logger.ErrorfCtx(c, "parse int form field, err: %v", err)
						continue
					}
					fieldValue.SetInt(int64(intVal))

				case reflect.Float32, reflect.Float64:
					f64, err := strconv.ParseFloat(formValue, 64)
					if err != nil {
						logger.ErrorfCtx(c, "parse float form field, err: %v", err)
						continue
					}
					fieldValue.SetFloat(f64)

				default:
					logger.ErrorfCtx(c, "unhandled default case")
				}

			}

		}

		callFuncName := ""

		lastFuncNames := strings.Split(runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name(), ".")
		if len(lastFuncNames) > 0 {
			callFuncName = lastFuncNames[len(lastFuncNames)-1]
		}

		if hasForm {
			logger.AccessfCtx(c, "=======> InPutLog: %v", callFuncName)
		} else {
			logger.AccessfCtx(c, "=======> InPutLog: %v, http body: %+v", callFuncName, in)
		}

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

			respJson, _ := sonic.Marshal(responseData)
			if !IsNoOutputLog(c) {
				logger.AccessfCtx(ctx, "<-------- %v, http response: %+v", callFuncName, string(respJson))
			}
			c.JSON(statusCode, responseData)
		}
		c.Next()
	}
}
