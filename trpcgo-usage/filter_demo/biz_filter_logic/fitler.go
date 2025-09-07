package bizfilterlogic

import (
	"context"
	"fmt"

	"trpc.group/trpc-go/tnet/log"
	"trpc.group/trpc-go/trpc-go/filter"
	trpcHttp "trpc.group/trpc-go/trpc-go/http"
	"trpc.group/trpc-go/trpc-go/plugin"
)

func init() {
	plugin.Register(pluginName, &pluginImp{})
}

// 自定义 filter 的 plugin
const (
	pluginName = "filterNameOne"
	pluginType = "filterTypeOne"
)

type pluginImp struct {
}

func (p *pluginImp) Type() string {
	return pluginType
}

type pluginConfig struct {
	ABC int    `yaml:"abc"`
	XYZ string `yaml:"xyz"`
}

func (p *pluginImp) Setup(name string, configDec plugin.Decoder) error {
	cfg := pluginConfig{}
	// log.Debugf("=======> decode cfg")
	fmt.Println("=======> decode cfg")
	if err := configDec.Decode(&cfg); err != nil {
		return err
	}
	// log.Debugf("=======> plugin filter config: %+v", cfg)
	fmt.Printf("=======> plugin filter, name: %v, config: %+v\n", name, cfg)
	filter.Register(name, serverFilterWrapper(&cfg), nil)

	return nil
}

func serverFilterWrapper(cfg *pluginConfig) filter.ServerFilter {
	return func(ctx context.Context, req interface{}, handler filter.ServerHandleFunc) (interface{}, error) {
		fmt.Println("...........................................")
		var head = trpcHttp.Head(ctx)
		// 非http请求
		if head == nil {
			fmt.Printf("req before config: %+v\n", *cfg)
			ret, err := handler(ctx, req)
			fmt.Printf("req after config: %+v\n", *cfg)
			return ret, err
		}
		// // 是否跳过OA验证(path白名单)
		// var r = head.Request
		// if r.URL != nil && o.isInExcludePath(r.URL.Path) {
		// 	return handler(ctx, req)
		// }
		// token, err := DefaultParseTokenFunc(ctx, req)
		// if err != nil {
		// 	return nil, err
		// }
		// customInfo, err := DefaultSigner.Verify(token)
		// if err != nil {
		// 	return nil, errs.NewFrameError(errs.RetServerAuthFail, err.Error())
		// }
		// // 将认证信息存储到ctx中用于业务侧使用
		// innerCtx := context.WithValue(ctx, AuthJwtCtxKey, customInfo)
		// return handler(innerCtx, req)
		log.Debugf("req before config: %+v", *cfg)
		ret, err := handler(ctx, req)
		log.Debugf("req after config: %+v", *cfg)
		return ret, err
	}
}
