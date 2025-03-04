package lib

import (
	"time"

	"server-transport-go-usage/lib/message"
)

// NodeConf 是 Node的配置信息
type NodeConf struct {
	// the endpoints with which this node will
	// communicate. Each endpoint contains zero or more channels
	Endpoints []EndpointConf

	// 消息类型集合(业务)
	Dialect []*message.CmdMessageInfo

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
