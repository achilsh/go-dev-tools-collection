package handler

import (
	logger "github.com/achilsh/go-dev-tools-collection/base-lib/log"
	"github.com/gin-gonic/gin"

	model "demo-service/model/http_model"
	lang "demo-service/service/common/language_def"
	"demo-service/service/utils/config"
)

const (
	WelcomeWordTypeFromCnf = 0
	WelcomeWordTypeFromDB  = 1
)

type WelcomeWordMessageHandler interface {
	GetMessage(lan string) []string
}

type WelcomeWordGetByConf struct {
}

func (c *WelcomeWordGetByConf) GetMessage(lan string) []string {
	return []string{}
}

type WelcomeWordGetByDB struct{}

func (c *WelcomeWordGetByDB) GetMessage(lan string) []string {
	if lan == "" {
		lan = "en"
	}

	wlang := config.GetGlobalConfig().WelWordLang
	if wlang == nil {
		logger.Errorf("not config welcome word language")
		return nil
	}
	//
	switch lan {
	case lang.LanguageKeyNameEn:
		return wlang.En
	case lang.LanguageKeyNameZh:
		return wlang.Zh
	default:
		return wlang.En
	}
	return []string{}
}

func GetWelcomeWordInstance(msgType int) WelcomeWordMessageHandler {
	var inst WelcomeWordMessageHandler
	switch msgType {
	case WelcomeWordTypeFromCnf:
		inst = new(WelcomeWordGetByConf)
	case WelcomeWordTypeFromDB:
		inst = new(WelcomeWordGetByDB)
	default:
		inst = new(WelcomeWordGetByConf)
	}
	return inst
}

// WelcomeWordHandle 欢迎词处理逻辑
func WelcomeWordHandle(ctx *gin.Context, _ *model.RequestWelcomeWordParam) (*model.ResponseWelcomeWordParam, error) {
	var (
		ret = &model.ResponseWelcomeWordParam{}
	)

	return ret, nil
}
