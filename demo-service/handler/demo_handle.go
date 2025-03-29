package handler

import (
	"fmt"

	logger "github.com/achilsh/go-dev-tools-collection/base-lib/log"
	"github.com/gin-gonic/gin"

	"demo-service/model/http_model"
	"demo-service/service/common/language_def"
)

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
