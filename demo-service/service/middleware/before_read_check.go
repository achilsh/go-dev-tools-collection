package middleware

import (
	"mime/multipart"

	logger "github.com/achilsh/go-dev-tools-collection/base-lib/log"
	"github.com/gin-gonic/gin"
)

func MarshalOpFormValue(dstMap map[string]any, src *multipart.Form) {
	// 解析普通表单字段
	for key, values := range src.Value {
		if len(values) == 1 {
			dstMap[key] = values[0]
		} else {
			dstMap[key] = values[0]
		}
	}
	logger.Debugf("data: %+v", dstMap)
}
func CheckBeforeReadBodyOnForm(ctx *gin.Context) (int, error) {
	logger.Debugf("call check read body form before.")
	// return 4001, fmt.Errorf("mock error: %v", 4001)
	return 0, nil
}

func CheckBeforeReadBodyOnJson(ctx *gin.Context) (int, error) {
	logger.Debugf("call check read body json before.")
	// return 4002, fmt.Errorf("mock error: %v", 4002)
	return 0, nil
}
