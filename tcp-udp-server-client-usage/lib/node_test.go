package lib

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type DemoMessageData struct {
	A int
	B float32
	C bool
}

func TestNodeDemo(t *testing.T) {
	var endPCnf []EndpointConf
	endPCnf = append(endPCnf, EndpointTCPServer{Address: ":5600"})
	//

	var endCnf = NodeConf{
		Endpoints:    endPCnf,
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 2 * time.Second,
		IdleTimeout:  10 * time.Second,
	}
	//
	endCnf.RegisterCmdMessage(1, new(DemoMessageData))
	endCnf.RegisterCmdMessage(2, new(DemoMessageData))

	n, err := NewNode(endCnf)
	assert.Nil(t, err)
	assert.NotNil(t, n)
	defer n.Close()
	select {
	case <-time.After(10 * time.Second):
		//
	}
}
