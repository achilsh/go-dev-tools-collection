package main

import (
	demo "server-transport-go-usage/gen/go/lib"
	"server-transport-go-usage/lib"
	"time"

	. "server-transport-go-usage/lib/utils"
)

func main() {
	LogPrintln("running server.....")
	var endPCnf []lib.EndpointConf
	endPCnf = append(endPCnf, lib.EndpointTCPServer{Address: "0.0.0.0:5600"})
	//
	var endCnf = lib.NodeConf{
		Endpoints:    endPCnf,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
		IdleTimeout:  0 * time.Second,
	}
	//
	lib.RegisterCmdMessage[demo.BizReqMsg](&endCnf, 1)
	lib.RegisterCmdMessage[demo.BizReqMsg](&endCnf, 2)

	n, err := lib.NewNode(endCnf)
	if err != nil {
		LogPrintln("new node fail, err: ", err)
		return
	}
	if n == nil {
		LogPrintln("new node is nil")
		return
	}
	defer n.Close()

	for evt := range n.Events() {
		if item, ok := evt.(*lib.EventFrame); ok {
			LogPrintf("msg seq: %v, msg Type: %v, msg data: %v\n",
				item.Frame.GetPkgSeq(), item.Frame.GetPkgType(), item.Frame.GetMessage())
			continue
		} else {
			if item, ok := evt.(*lib.EventParseError); ok {
				LogPrintf("receive err parse message: %v", item)
				continue
			}
			if item, ok := evt.(*lib.EventChannelOpen); ok {
				LogPrintf("receive open new connect event, %v", item)
				continue
			}
		}
		LogPrintf("receive msg, value: %v", evt)
	}
}
