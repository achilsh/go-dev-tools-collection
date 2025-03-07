package main

import (
	"time"


	demo "server-transport-go-usage/gen/go/proto"
	"server-transport-go-usage/lib"
	"server-transport-go-usage/lib/message"
	. "server-transport-go-usage/lib/utils"
)

func main() {
	var testDemoSeq uint16 = 0
	LogPrintln("running client .....")
	var endPCnf []lib.EndpointConf
	endPCnf = append(endPCnf, lib.EndpointUDPClient{Address: "0.0.0.0:5601"})
	var endCnf = lib.NodeConf{
		Endpoints:    endPCnf,
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 2 * time.Second,
		IdleTimeout:  2 * time.Second,
	}
	lib.RegisterCmdMessage[demo.BizReqMsg](&endCnf, 1)
	lib.RegisterCmdMessage[demo.BizReqMsg](&endCnf, 2)

	n, err := lib.NewNode(&endCnf)
	if err != nil {
		LogPrintln("new client end node fail, err: ", err)
		return
	}
	if n == nil {
		LogPrintln("new client end node is nil.")
		return
	}
	defer n.Close()

	mockInitCost := func() {
		time.Sleep(100 * time.Millisecond)
	}
	mockInitCost()

	for i := 0; i < 3; i++ {
		testDemoSeq++
		toSendMsg := &message.DecodedMessage{
			HeaderMessage: &message.HeaderMessage{
				StartFlag: message.PKG_START_FLAG,
				PkgSeq:    testDemoSeq,
				DevType:   1,
				PkgType:   1,
			},
			DecodedMsg: &demo.BizReqMsg{
				Name: "achilsh",
				Age:  120,
			},
		}

		err = n.WriteFrameAll(toSendMsg)
		if err != nil {
			LogPrintln("write frame data fail, err: ", err)
		}
	}
	//空闲一段时间。
	time.Sleep(6*time.Second)
	for i := 0; i < 2; i++ {
		testDemoSeq++
		toSendMsg := &message.DecodedMessage{
			HeaderMessage: &message.HeaderMessage{
				StartFlag: message.PKG_START_FLAG,
				PkgSeq:    testDemoSeq,
				DevType:   1,
				PkgType:   1,
			},
			DecodedMsg: &demo.BizReqMsg{
				Name: "achilsh",
				Age:  120,
			},
		}

		err = n.WriteFrameAll(toSendMsg)
		if err != nil {
			LogPrintln("write frame data fail, err: ", err)
		}
	}
	//等待一段时间客户端退出 
	time.Sleep(4*time.Second)
	LogPrintf("client exit now.")
}
