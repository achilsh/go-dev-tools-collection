package rewrites

import (
	"fmt"
	balanceupstreams "http_gateway_demo/balance_upstreams"
	"regexp"

	"github.com/zeromicro/go-zero/gateway"
)

type RewriteHttpPathCfgElem struct {
	Method string `json:"Method"`
	From   string `json:"From"`
	To     string `json:"To"`
}

type RewritHttpUrlPathConf struct {
	ReWrites []RewriteHttpPathCfgElem `json:"ReWrites,optional"`
	// 规则定义
	RewriteRules []RewriteRule `json:"-"`
}

type ServiceConfig struct {
	gateway.GatewayConf
	RewritHttpUrlPathConf
	balanceupstreams.UpstreamsBizConfig
}

type RewriteRule struct {
	From   *regexp.Regexp // from是正则表达式
	To     string         // to 是目标path
	Method string
}

func TransRewriteCfgToRules(rewriteCfg *RewritHttpUrlPathConf) *RewritHttpUrlPathConf {
	if rewriteCfg == nil {
		return nil
	}
	rules := make([]RewriteRule, 0, len(rewriteCfg.ReWrites))
	for _, rcfg := range rewriteCfg.ReWrites {
		re, err := regexp.Compile(rcfg.From)
		if err != nil {
			panic(fmt.Sprintf("compile from rege fail, from: %v", rcfg.From))
		}
		rules = append(rules, RewriteRule{
			From:   re,
			To:     rcfg.To,
			Method: rcfg.Method,
		})
	}
	rewriteCfg.RewriteRules = rules
	return rewriteCfg
}
