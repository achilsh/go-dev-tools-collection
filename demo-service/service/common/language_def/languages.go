package language_def

import "github.com/gin-gonic/gin"

type LanguageDescType map[string]string

const XLanguageHeaderKey = "Language"
const (
	LanguageKeyNameEn = "en"
	LanguageKeyNameDe = "de"
	LanguageKeyNameEs = "es"
	LanguageKeyNameFr = "fr"

	LanguageKeyNameIt = "it"

	LanguageKeyNameJa = "ja"

	LanguageKeyNameZh = "zh"
)

var LanguageDesc LanguageDescType = map[string]string{
	LanguageKeyNameEn: "english",
	LanguageKeyNameDe: "de",
	LanguageKeyNameEs: "es",
	LanguageKeyNameFr: "fr",
	LanguageKeyNameIt: "it",
	LanguageKeyNameJa: "ja",
	LanguageKeyNameZh: "homeland_china",
}

func GetLanguage(ctx *gin.Context) string {
	v, exist := ctx.Get(XLanguageHeaderKey)
	if exist == false {
		return "zh"
	}
	return v.(string)
}
