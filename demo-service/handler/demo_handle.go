package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	logger "github.com/achilsh/go-dev-tools-collection/base-lib/log"
	"github.com/gin-gonic/gin"

	"github.com/achilsh/go-dev-tools-collection/demo-service/model/http_model"
	"github.com/achilsh/go-dev-tools-collection/demo-service/service/common/language_def"
	rw "github.com/achilsh/go-dev-tools-collection/gin-usage-wrapper/router_wrapper"
)

// form 数据上传处理：

func DemoFormInOut(ctx *gin.Context, in *http_model.FormReqParam) (*http_model.FormResponse, error) {
	for fileName := range in.FileContentMap {
		logger.Infof("file name: %v", fileName)
	}

	responseContent := ""
	for fileName, fileContent := range in.FileContentMap {
		if len(fileContent) > 0 {
			logger.Infof("receive clent request data len: %v", len(fileContent))

			f, err := os.Create(fileName)
			if err != nil {
				logger.Errorf("create file fail, err: %v", err)
				return nil, fmt.Errorf("create file fail")
			}
			defer f.Close()

			readerBuf := bytes.NewBuffer(fileContent)
			_, err = io.Copy(f, readerBuf)
			if err != nil {
				logger.Errorf("write recv buf to file fail, err: %v", err)
				return nil, fmt.Errorf("write recv data to fail.")
			}
			if len(responseContent) > 0 {
				responseContent += ", " + fmt.Sprintf("fileName: %v, contentLen: %v", fileName, len(fileContent))
			} else {
				responseContent += fmt.Sprintf("fileName: %v, contentLen: %v", fileName, len(fileContent))
			}
		}
	}
	return &http_model.FormResponse{
		ABC: 0,
		XYZ: responseContent,
	}, nil

	// return nil, fmt.Errorf("not processed....")

}
func DemoIn(ctx *gin.Context) (int, error) {
	l := language_def.GetLanguage(ctx)
	logger.Infof("language: %v", l)
	return 1000, nil
}

func DemoInPtr(ctx *gin.Context) (*int, error) {
	l := language_def.GetLanguage(ctx)
	logger.Infof("language: %v", l)
	x := 1000

	_ = x
	return nil, fmt.Errorf("this is error: %v", x)
}

type StreamData struct {
	A int    `json:"aa"`
	B string `json:"bb"`
}

type DemoZRequest struct {
	AA  int `json:"aa"`
	Age int `json:"age"`
}
type DemoZResponse struct {
	Address string `json:"address"`
}

func DemoABC(ctx *gin.Context, in *DemoZRequest) (*DemoZResponse, error) {
	rw.FlagNoOutputLog(ctx)

	ret := &DemoZResponse{
		Address: "sfadfadf------",
	}

	logger.Infof("xxxx--->: %+v", ret)
	// return ret, nil
	return nil, fmt.Errorf("define error, 2000")
}

func DemoInOut(ctx *gin.Context, in *http_model.RequestParam) (*http_model.ResponseParam, error) {
	ret := &http_model.ResponseParam{
		Result: fmt.Sprintf("id: %v, name: %v", in.Id, in.Name),
	}

	// 下面增加的是流式处理， 测试用例：  curl -N -X POST http://localhost:5657/demo-server/v1/x2   -H "Content-Type: application/json"   -d '{"name":"你好"}'
	f, ok := ctx.Writer.(http.Flusher)
	if !ok {
		logger.Errorf("not support flusher.")
	} else {
		logger.Info("support flusher.")
	}
	logger.Accessf("this is access log......., time: %+v", time.Now())

	ctx.Header("Content-Type", "text/event-stream")
	ctx.Header("Cache-Control", "no-cache")
	ctx.Header("Connection", "keep-alive")
	// ctx.Header("Access-Control-Allow-Origin", "*") // 允许跨域访问
	//
OUTEXIT:
	for {
		select {
		case <-ctx.Done():
			break OUTEXIT
		default:
			for i := 0; i < 10; i++ {
				var abc = StreamData{
					A: i,
					B: fmt.Sprintf("call: %v", i),
				}
				c, _ := json.Marshal(&abc)
				fmt.Fprintf(ctx.Writer, "xxxx---xxx...: %v", string(c))
				i++
				f.Flush()
				logger.Info("flush data to client.")
				logger.Accessf("this is access log......., time: %+v", time.Now())
				time.Sleep(10 * time.Millisecond)
				logger.InfofCtx(ctx, "----------------: %v", 1)
				logger.DebugfCtx(ctx, "--------- debug ----------: %v", 2)
				logger.ErrorfCtx(ctx, "--------- error ----------: %v", 3)
			}
			break OUTEXIT
		}
	}
	return ret, nil
}

func DemoNoOut(ctx *gin.Context, in *http_model.RequestParam) error {
	logger.Infof("in data: %+v", *in)
	logger.InfofCtx(ctx, "with context: is data: %+v", *in)
	return fmt.Errorf("no out error: 3000")
}

func DemoGetInOut(ctx *gin.Context, in *http_model.RequestParam) (*http_model.ResponseParam, error) {
	ret := &http_model.ResponseParam{
		Result: fmt.Sprintf("id: %v, name: %v", in.Id, in.Name),
	}
	return ret, nil
}

func DemoGetNoIn(ctx *gin.Context) (int, error) {
	l := language_def.GetLanguage(ctx)
	logger.Infof("language: %v", l)
	return 1000, nil
}
func DemoGetNoInPtr(ctx *gin.Context) (*int, error) {
	l := language_def.GetLanguage(ctx)
	logger.Infof("language: %v", l)
	x := 1000
	return &x, nil
}

func DemoGetNoOut(ctx *gin.Context, in *http_model.RequestParam) error {
	logger.Infof("in data: %+v", *in)
	return nil
}
