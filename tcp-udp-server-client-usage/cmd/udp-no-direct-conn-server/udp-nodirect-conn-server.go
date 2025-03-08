package main

import (
	"time"

	demo "server-transport-go-usage/gen/go/proto"
	"server-transport-go-usage/lib"
	. "server-transport-go-usage/lib/utils"
)

func main() {
	LogPrintln("running udp server.....")
	var endPCnf []lib.EndpointConf
	endPCnf = append(endPCnf, lib.EndPointUDPNoDirectServer{Address: "0.0.0.0:5603"})
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

outStep:
	// 需要注册数据包的业务逻辑：
	// BizHandle(context.Context, *EventFrame) (error, any)
	for evt := range n.Events() {
		switch dataFrame := evt.(type) {
		case *lib.EventFrame:
			LogPrintf("msg seq: %v, msg Type: %v, msg data: %v\n",
				dataFrame.Frame.GetPkgSeq(), dataFrame.Frame.GetPkgType(), dataFrame.Frame.GetMessage())
		case *lib.EventParseError:
			LogPrintf("receive err parse message: %v", dataFrame)

		case *lib.EventChannelOpen:
			LogPrintf("receive open new connect event, %v", dataFrame)

		case *lib.EventChannelClose:
			LogPrintf("receive close connect event, %v", dataFrame)
			goto outStep
		default:
			LogPrintf("receive msg, value: %v", evt)
			goto outStep
		}
	}
	LogPrintf("exist on server mian.")
}