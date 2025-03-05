package lib

import (
	"reflect"
	"server-transport-go-usage/lib/message"
	"time"

	. "server-transport-go-usage/lib/utils"
)

func RegisterCmdMessage[M any](ncf *NodeConf, cmd uint16) {
	ncf.RegisterCmdMessage(cmd, new(M))
}

// NodeConf 是 Node的配置信息
type NodeConf struct {
	// the endpoints with which this node will
	// communicate. Each endpoint contains zero or more channels
	Endpoints []EndpointConf

	// 消息类型集合(业务); 消息的类型要保持全局唯一
	Dialect map[uint16]*message.CmdMessageInfo

	// (optional) read timeout.
	// It defaults to 10 seconds.
	ReadTimeout time.Duration
	// (optional) write timeout.
	// It defaults to 10 seconds.
	WriteTimeout time.Duration
	// (optional) timeout before closing idle connections.
	// It defaults to 60 seconds.
	IdleTimeout time.Duration
}

// RegisterCmdMessage add message struct format on this cmd.
func (nf *NodeConf) RegisterCmdMessage(cmd uint16, msg any) {
	if nf.Dialect == nil {
		nf.Dialect = make(map[uint16]*message.CmdMessageInfo)
	}
	if reflect.TypeOf(msg).Kind() != reflect.Pointer {
		LogPrintf("cmd: %v, msg kind: %v is not pointer\n", cmd, reflect.TypeOf(msg).Kind())
		return
	}

	if _, ok := nf.Dialect[cmd]; ok {
		return
	}

	nf.Dialect[cmd] = &message.CmdMessageInfo{
		Cmd:     cmd,
		MsgType: reflect.TypeOf(msg),
	}
}
