package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type (
	ContextValueKey string
)

const (
	CtxRequestID      = "X-Request-Id"
	CtxReqBody        = "req_body"
	CtxUserID         = "ctx_user_id"
	CtxLoginUser      = "ctx_login_user"
	CtxSOPCountry     = "ctx_sop_country"
	CtxSOPAuthCodes   = "ctx_sop_auth_codes"
	CtxSOPAllAuth     = "ctx_sop_all_auth"
	CtxBDCenterRegion = "bd-region"
	CtxStatus         = "ctx_status"
)

func GetRequestID(ctx *gin.Context) string {
	id, exist := ctx.Get(CtxRequestID)
	if !exist {
		return ""
	}
	return id.(string)
}

func GenerateRequestID() string {
	id := uuid.New()
	s := id.String()
	return s
}

func SetRequestBody(ctx *gin.Context, body []byte) {
	ctx.Set(CtxReqBody, body)
}
func GetRequestBody(ctx *gin.Context) []byte {
	data, exist := ctx.Get(CtxReqBody)
	if !exist {
		return nil
	}
	return data.([]byte)
}
