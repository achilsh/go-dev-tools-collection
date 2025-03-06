package lib

import (
	"context"
	"fmt"

	. "server-transport-go-usage/lib/utils"
)

type LogicProcesser interface {
	Handle(ctx context.Context, in Event) (error, any)
}
type BizLogicRouter struct {
	//key: 
	routerProcess map[string] LogicProcesser
}

func NewBizLogicRouter() *BizLogicRouter{
	return &BizLogicRouter{
		routerProcess: make(map[string]LogicProcesser),
	}
}


// AddRouter 需要手动添加路由处理函数
func (br *BizLogicRouter) AddRouter(cmd string, handle LogicProcesser) {
	if _, ok := br.routerProcess[cmd]; !ok {
		 br.routerProcess[cmd]= handle
	} else {
		LogPrintf("cmd %v has register handle.", cmd)
	}
}
func (br *BizLogicRouter) GetRouterHandle(cmd string)(LogicProcesser, error) {
	if cmd == "" {
		return nil, fmt.Errorf("cmd is empty")
	}
	if br ==nil || br.routerProcess == nil {
		return  nil, fmt.Errorf("router manager not init")
	}

	item, ok := br.routerProcess[cmd]
	if !ok {
		return nil, fmt.Errorf("not register router handle for cmd: %v", cmd)
	}
	return  item, nil
}