package middleware

import (
	"github.com/gin-gonic/gin"

	"github.com/achilsh/go-dev-tools-collection/demo-service/service/common/language_def"
)

func LanguageManger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		lang := ctx.Request.Header.Get(language_def.XLanguageHeaderKey)
		if lang == "" {
			lang = language_def.LanguageKeyNameZh
		}
		ctx.Set(language_def.XLanguageHeaderKey, lang)
		ctx.Next()
	}
}
