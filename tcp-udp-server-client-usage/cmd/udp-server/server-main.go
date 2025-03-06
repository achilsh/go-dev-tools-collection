package main 

import (
	demo "server-transport-go-usage/gen/go/proto"
	"server-transport-go-usage/lib"
	"time"

	. "server-transport-go-usage/lib/utils"
)

func main() {
	LogPrintln("running udp server.....")
	var endPCnf []lib.EndpointConf
	endPCnf = append(endPCnf, lib.EndpointUDPServer{Address: "0.0.0.0:5601"})
	//
	var endCnf = lib.NodeConf{
		Endpoints:    endPCnf,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  5 * time.Second,
	}
	//
	lib.RegisterCmdMessage[demo.BizReqMsg](&endCnf, 1)
	lib.RegisterCmdMessage[demo.BizReqMsg](&endCnf, 2)

	n, err := lib.NewNode(&endCnf)
	if err != nil {
		LogPrintln("new node fail, err: ", err)
		return
	}
	if n == nil {
		LogPrintln("new node is nil")
		return
	}
	defer n.Close()
	defer func() {
		if e := recover(); e != nil {
			LogPrintln("panic: ", e)
		}
	}()

	// 需要注册数据包的业务逻辑：
	// BizHandle(context.Context, *EventFrame) (error, any)
	for evt := range n.Events() {
		if item, ok := evt.(*lib.EventFrame); ok {
			LogPrintf("msg seq: %v, msg Type: %v, msg data: %v\n",
				item.Frame.GetPkgSeq(), item.Frame.GetPkgType(), item.Frame.GetMessage())
			continue
		}
		if item, ok := evt.(*lib.EventParseError); ok {
			LogPrintf("receive err parse message: %v", item)
			continue
		}
		if item, ok := evt.(*lib.EventChannelOpen); ok {
			LogPrintf("receive open new connect event, %v", item)
			continue
		}

		LogPrintf("receive msg, value: %v", evt)
	}
	LogPrintf("exist on server mian.")
}