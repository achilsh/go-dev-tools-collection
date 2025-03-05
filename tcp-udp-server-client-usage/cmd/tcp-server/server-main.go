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
	endCnf.RegisterCmdMessage(1, new(demo.BizReqMsg))
	endCnf.RegisterCmdMessage(2, new(demo.BizReqMsg))

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
		} else {
			if item, ok := evt.(*lib.EventParseError); ok {
				LogPrintf("receive err parse message: %v", item)
			}
			LogPrintln("receive msg.")
		}

	}
}
