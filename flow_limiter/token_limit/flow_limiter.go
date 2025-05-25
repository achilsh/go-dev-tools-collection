package tokenlimit

import (
	"context"
	"embed"
	"fmt"
	"time"

	logger "github.com/achilsh/go-dev-tools-collection/base-lib/log"
	redisLib "github.com/achilsh/go-dev-tools-collection/redis-wrapper/lib"
)

//go:embed token_limiter.lua
var tokenLimiterLuaFS embed.FS 
var tokenLimiterLuaContent []byte


func init() {
	var err error 
	tokenLimiterLuaContent, err = tokenLimiterLuaFS.ReadFile("token_limiter.lua")
	if err !=nil {
		fmt.Printf("load token_limiter.lua fail, err: %v\n", err)
		return 
	}
	fmt.Printf("load token_limiter.lua content: %v\n", string(tokenLimiterLuaContent))
}


type TokenLimiterParam struct {
	 maxTokenNums int  //每秒钟请求最大个数
	 refillRateInSecond int //单位为秒时间段； 每一秒中添加token的数量
	 maxRequestNums int  //每日请求个数，按00-24时刻算
	 KeyRateControl string  //速率控制的key； 依赖于业务输入
	 KeyTotalRequest string //总数控制的key. 依赖于业务输入
}

type FillParamsFunc func(*TokenLimiterParam)

func InitKeys(rateKey, totalReqKey string) FillParamsFunc {
	return func(item *TokenLimiterParam) {
		if item == nil{
			return 
		}
		//
		item.KeyRateControl = rateKey 
		item.KeyTotalRequest = totalReqKey
	}
}

// 
func InitControlInfo(maxTokens  int, refill int, maxDaily  int) FillParamsFunc {
	return func(item *TokenLimiterParam) {
		if item ==nil {
			return 
		}

		item.maxTokenNums = maxTokens  
		item.maxRequestNums =  maxDaily  
		item.refillRateInSecond = refill
	}
}

type TokenLimiter struct {
	 // 这两个是基本固定下来
	hashLua string 
	redisClient redisLib.RedisOps
}


type  InitTokenLimiter func(*TokenLimiter)

func InitRedis(rdsOps ...redisLib.RedisOption) InitTokenLimiter {
	return func(item *TokenLimiter) {
		if item==nil {
			return 
		}
		item.redisClient = redisLib.NewRedisClient(rdsOps...)
	}
}


func LoadHashLua() InitTokenLimiter {
	return func(item *TokenLimiter) {
		if item == nil {
			return 
		}

		redisHandle := item.redisClient.GetRedisRow()
		if redisHandle == nil {
			logger.Errorf("redis not init")
			return 
		}
		
		hashStr, err :=redisHandle.ScriptLoad(context.Background(), string(tokenLimiterLuaContent)).Result()
		if err != nil{
			logger.Errorf("calc hash for lua fail, err: %v", err)
			return 
		}
		item.hashLua = hashStr
	}
}

// need 
func NewTokenLimiter(opts ...InitTokenLimiter) *TokenLimiter {
	tler := new(TokenLimiter)
	for _, op := range opts {
		op(tler)
	}
	return tler
}



// how to use? 
func (tl *TokenLimiter)Allow(params ...FillParamsFunc) bool {
	rdsHandle := tl.redisClient.GetRedisRow()
	if rdsHandle == nil {
		logger.Errorf("get redis handle is nil, maybe not init redis.")
		return false
	}

	var  oneParam TokenLimiterParam = TokenLimiterParam{
	maxTokenNums: 10,  //每秒钟请求最大个数
	 refillRateInSecond : 1,  //单位为秒时间段
	 maxRequestNums :20 , //每日请求个数，按00-24时刻算
	//  KeyRateControl string  //速率控制的key； 依赖于业务输入
	//  KeyTotalRequest string //总数控制的key. 依赖于业务输入
	}

	for _, pFunc := range params {
		pFunc(&oneParam)
	}
	keys := []string{oneParam.KeyRateControl, oneParam.KeyTotalRequest}
	args :=[]any{oneParam.maxTokenNums, oneParam.refillRateInSecond, oneParam.maxRequestNums, time.Now().Unix()}

	results, err := rdsHandle.EvalSha(context.Background(), tl.hashLua, keys, args ).Result()
	if err != nil {
		return false
	}

	result := results.(int64)
	if result < 0 {
		logger.Errorf("check result: %v", result)
		return false
	}

	return true
}
// func NewToken






