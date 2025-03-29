package middleware

import (
	"fmt"
	"math/big"

	"github.com/gin-gonic/gin"
)

func BcCompSafe(num1, num2 string, scale *int) (int, error) {
	a := new(big.Rat)
	if _, ok := a.SetString(num1); !ok {
		return 0, fmt.Errorf("无效的数字: %s", num1)
	}

	b := new(big.Rat)
	if _, ok := b.SetString(num2); !ok {
		return 0, fmt.Errorf("无效的数字: %s", num2)
	}

	if scale != nil {
		prec := *scale
		aStr := a.FloatString(prec)
		bStr := b.FloatString(prec)

		if _, ok := a.SetString(aStr); !ok {
			return 0, fmt.Errorf("格式化数字失败: %s", aStr)
		}
		if _, ok := b.SetString(bStr); !ok {
			return 0, fmt.Errorf("格式化数字失败: %s", bStr)
		}
	}

	return a.Cmp(b), nil
}

func CheckTokenTimeStamp(ctx *gin.Context) string {

	return ""
}

func getAccessTokenRedisKey(uid string) string {
	return fmt.Sprintf("%v%v", "access:tk:", uid)
}

func CheckUidCache(ctx *gin.Context) string {
	payload := GetPayLoadFromCtx(ctx)
	if payload == nil {
		return "4029"
	}
	if payload.Akey == "" || payload.Uid == "" {
		return "4029"
	}
	return ""
}
