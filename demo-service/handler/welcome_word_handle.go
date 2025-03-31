package handler

import (
	logger "github.com/achilsh/go-dev-tools-collection/base-lib/log"
	"github.com/gin-gonic/gin"

	model "demo-service/model/http_model"
	lang "demo-service/service/common/language_def"
	"demo-service/service/error_code"
	"demo-service/service/utils/config"
	"demo-service/service/utils/routinues"
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
		wg  = routinues.NewRoutineGroupWrap()
		ret = &model.ResponseWelcomeWordParam{}

		errGetWord error = nil
	)
	lanVal := lang.GetLanguage(ctx)

	wg.AsyncTimeoutRun(true, 500, func() {
		// wg.AsyncRun(true, func() {
		welcomeWord := GetWelcomeWordInstance(WelcomeWordTypeFromCnf).GetMessage(lanVal)
		if len(welcomeWord) == 0 {
			errGetWord = error_code.WelcomeWordError
			logger.Warnf("not get valid welcome word for lang: %v", lanVal)
			return
		}
		ret.WelcomeWords = welcomeWord
	})

	wg.AsyncTimeoutRun(true, 500, func() {
		handle := &GuessQuestionFromDB{}
		guessQ := handle.GetGuessQuestions(lanVal)
		if len(guessQ) == 0 {
			logger.Warnf("get guess question is empty")
			return
		}
		ret.ToGuessQuestions = guessQ
	})
	//
	wg.Wait()
	if errGetWord != nil {
		return nil, errGetWord
	}

	return ret, nil
}
