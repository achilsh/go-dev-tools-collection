package lib

import (
	"fmt"
	demo "server-transport-go-usage/gen/go/proto"
	"server-transport-go-usage/lib/message"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	. "server-transport-go-usage/lib/utils"
)

var testDemoSeq uint16 = 0

func TestTcpClient(t *testing.T) {
	var endPCnf []EndpointConf
	endPCnf = append(endPCnf, EndpointTCPClient{Address: "0.0.0.0:5600"})
	var endCnf = NodeConf{
		Endpoints:    endPCnf,
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 2 * time.Second,
		IdleTimeout:  10 * time.Second,
	}
	endCnf.RegisterCmdMessage(1, new(demo.BizReqMsg))
	endCnf.RegisterCmdMessage(2, new(demo.BizReqMsg))
	n, err := NewNode(&endCnf)
	assert.Nil(t, err)
	assert.NotNil(t, n)
	defer n.Close()
	//
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

	time.Sleep(3 * time.Second)
	err = n.WriteFrameAll(toSendMsg)
	assert.Nil(t, err)
}
func TestNodeServer(t *testing.T) {
	fmt.Println("aaa.....")
	LogPrintln("running server.....")
	var endPCnf []EndpointConf
	endPCnf = append(endPCnf, EndpointTCPServer{Address: "0.0.0.0:5600"})
	//
	var endCnf = NodeConf{
		Endpoints:    endPCnf,
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 2 * time.Second,
		IdleTimeout:  10 * time.Second,
	}
	//
	endCnf.RegisterCmdMessage(1, new(demo.BizReqMsg))
	endCnf.RegisterCmdMessage(2, new(demo.BizReqMsg))

	n, err := NewNode(&endCnf)
	assert.Nil(t, err)
	assert.NotNil(t, n)
	defer n.Close()

	for evt := range n.Events() {
		if item, ok := evt.(*EventFrame); ok {
			LogPrintf("msg type: %v, msg data: %v\n", item.Frame.GetDevType(), item.Frame.GetMessage())
		}
		LogPrintln("receive msg.")
	}
}
