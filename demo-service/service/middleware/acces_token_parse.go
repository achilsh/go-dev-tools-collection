package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type AccessTokenPayload struct {
	Exp    int    `json:"exp"`  //过期时间
	SubStr string `json:"sub"`  //信息主题
	Akey   string `json:"akey"` //token标识
	AoFn   int    `json:"aofn"` //aes key iv 的偏移量
	JoFn   int    `json:"jofn"` //JWT密钥偏移量
	Uid    string `json:"uid"`  //加密后的 uid
	Iss    string `json:"iss"`  //签发方
	Nbf    int    `json:"nbf"`  //当前时间 生效时间
	Iat    int    `json:"iat"`  //Issued at，签发时间
}

const (
	BodyPayloadKeyName = "body_payload"
	BodyUidKeyName     = "uid_value"
)

func SetPayLoadToCtx(payLoad *AccessTokenPayload, ctx *gin.Context) {
	ctx.Set(BodyPayloadKeyName, payLoad)
	ctx.Set(BodyUidKeyName, payLoad.Uid)
}
func GetPayLoadFromCtx(ctx *gin.Context) *AccessTokenPayload {
	if ctx == nil {
		return nil
	}
	payload, ok := ctx.Get(BodyPayloadKeyName)
	if !ok {
		return nil
	}
	ret := payload.(*AccessTokenPayload)
	if ret == nil {
		return nil
	}
	return ret
}

func ParseToken(bodyMap map[string]any, ctx *gin.Context) (string, *AccessTokenPayload) {
	accessToken, err := CheckAccessToken(bodyMap)
	if err != nil || accessToken == "" {
		return "4001", nil
	}

	var payloadDecode AccessTokenPayload

	return "", &payloadDecode
}
func CheckAccessToken(bodyMap map[string]any) (string, error) {
	accessToken, ok := bodyMap[AccessTOKEN_FIELD_NAME]
	if !ok {
		return "", fmt.Errorf("access token is nil")
	}
	return accessToken.(string), nil
}
