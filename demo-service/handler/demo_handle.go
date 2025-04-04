package handler

import (
	"bytes"
	"fmt"
	"io"
	"os"

	logger "github.com/achilsh/go-dev-tools-collection/base-lib/log"
	"github.com/gin-gonic/gin"

	"demo-service/model/http_model"
	"demo-service/service/common/language_def"
)

// form 数据上传处理：

func DemoFormInOut(ctx *gin.Context, in *http_model.FormReqParam) (*http_model.FormResponse, error) {
	logger.Infof("file name: %v", in.FileName)
	if len(in.FileContent) > 0 {
		logger.Infof("receive clent request data len: %v", len(in.FileContent))

		f, err := os.Create(in.FileName)
		if err != nil {
			logger.Errorf("create file fail, err: %v", err)
			return nil, fmt.Errorf("create file fail")
		}
		defer f.Close()

		readerBuf := bytes.NewBuffer(in.FileContent)
		wlen, err := io.Copy(f, readerBuf)
		if err != nil {
			logger.Errorf("write recv buf to file fail, err: %v", err)
			return nil, fmt.Errorf("write recv data to fail.")
		}
		return &http_model.FormResponse{
			ABC: int(wlen),
		}, nil
	}
	return nil, fmt.Errorf("not processed....")

}
func DemoIn(ctx *gin.Context) (int, error) {
	l := language_def.GetLanguage(ctx)
	logger.Infof("language: %v", l)
	return 1000, nil

}

func DemoInOut(ctx *gin.Context, in *http_model.RequestParam) (*http_model.ResponseParam, error) {
	ret := &http_model.ResponseParam{
		Result: fmt.Sprintf("id: %v, name: %v", in.Id, in.Name),
	}
	return ret, nil
}

func DemoNoOut(ctx *gin.Context, in *http_model.RequestParam) error {
	logger.Infof("in data: %+v", *in)
	return nil
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

func DemoGetNoOut(ctx *gin.Context, in *http_model.RequestParam) error {
	logger.Infof("in data: %+v", *in)
	return nil
}
