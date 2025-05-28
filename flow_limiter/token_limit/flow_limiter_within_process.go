package tokenlimit

import (
	"golang.org/x/time/rate"
)

//  使用golang 标准库来限制单个进程内的并发数

type FlowLimiterInProcessParam struct {
	TokenPerSecond int //每秒生成的token数，
	MaxBurstTokens int //突发流量最大数，也可表示令牌桶的最多可以存储的令牌数， 是瞬间突发请求最大数。
}

type InitFlowLimiterProcessFunc func(*FlowLimiterInProcessParam)

func WithTokens(tokenNumsPerSecond int, tokenMaxBurstNums int) InitFlowLimiterProcessFunc {
	return func(item *FlowLimiterInProcessParam) {
		if item == nil {
			return
		}
		item.TokenPerSecond = tokenNumsPerSecond
		item.MaxBurstTokens = tokenMaxBurstNums
	}
}

type FlowLimiterOneProcess struct {
	limiter *rate.Limiter
}

func NewProcessFLowLimiter(opts ...InitFlowLimiterProcessFunc) *FlowLimiterOneProcess{
	params :=&FlowLimiterInProcessParam{}
	for _, opt := range opts {
		opt(params)
	}

	lm := rate.NewLimiter(rate.Limit(params.TokenPerSecond),params.MaxBurstTokens)
	return &FlowLimiterOneProcess{
		limiter:  lm,
	}
}


func (flp *FlowLimiterOneProcess) Allow() bool {
	if flp.limiter == nil {
		return false
	}
	return flp.limiter.Allow()
}