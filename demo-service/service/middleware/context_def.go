package middleware

import "github.com/gin-gonic/gin"

type (
	ContextValueKey string
)

const (
	CtxReqBody        = "req_body"
	CtxUserID         = "ctx_user_id"
	CtxLoginUser      = "ctx_login_user"
	CtxSOPCountry     = "ctx_sop_country"
	CtxSOPAuthCodes   = "ctx_sop_auth_codes"
	CtxSOPAllAuth     = "ctx_sop_all_auth"
	CtxBDCenterRegion = "bd-region"
	CtxStatus         = "ctx_status"
)

func SetRequestBody(ctx *gin.Context, body []byte) {
	ctx.Set(CtxReqBody, body)
}
func GetRequestBody(ctx *gin.Context) []byte {
	data, exist := ctx.Get(CtxReqBody)
	if exist == false {
		return nil
	}
	return data.([]byte)
}
