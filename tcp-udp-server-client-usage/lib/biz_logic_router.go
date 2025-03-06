package lib

import (
	"context"
	"fmt"

	. "server-transport-go-usage/lib/utils"
)

func GetProcessId(id uint16) string {
	if  id <= 0 {
		return ""
	}
	return fmt.Sprintf("pKey:%v", id)
}
// 需要业务编写实例： 所有的协议类型都要实现 处理该消息的逻辑
type LogicProcesser interface {
	Handle(ctx context.Context, in Event) (error, any)
}

// 需要业务编写实例：只要统一实现 解析失败，连接成功的处理函数
type ProcessWithBizInfoer interface {
	HandleParseFail(ctx context.Context, in Event) (error, any)
	HandleNewConn(ctx context.Context, in Event) (error, any)
}
type BizLogicRouter struct {
	//key: 
	routerProcess map[string] LogicProcesser
	withBizInfoProcess ProcessWithBizInfoer
}

// NewBizLogicRouter 需要业务初始化调用。
func NewBizLogicRouter() *BizLogicRouter{
	return &BizLogicRouter{
		routerProcess: make(map[string]LogicProcesser),
	}
}


// AddRouter 需要业务 初始化调用： 使用 cmd 和 上面 实例化LogicProcesser的对象一起，注入到流程中。
func (br *BizLogicRouter) AddRouter(cmd string, handle LogicProcesser) {
	if _, ok := br.routerProcess[cmd]; !ok {
		 br.routerProcess[cmd]= handle
	} else {
		LogPrintf("cmd %v has register handle.", cmd)
	}
}


// AddWithoutBizRouter 需要业务初始化调用： 使用 上面 实例化 ProcessWithBizInfoer 的对象一起，注入到流程中。
func (br *BizLogicRouter)AddWithoutBizRouter(handle ProcessWithBizInfoer) {
	br.withBizInfoProcess = handle
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
func (br* BizLogicRouter) GetBizNoHandle() ProcessWithBizInfoer {
	return br.withBizInfoProcess
}