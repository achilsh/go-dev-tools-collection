package lib

import (
	"time"

	"server-transport-go-usage/lib/message"
)

// NodeConf allows to configure a Node.
type NodeConf struct {
	// the endpoints with which this node will
	// communicate. Each endpoint contains zero or more channels
	Endpoints []EndpointConf

	// (optional) the dialect which contains the messages that will be encoded and decoded.
	// If not provided, messages are decoded in the MessageRaw struct.
	Dialect []*message.CmdMessageInfo

	// // (optional) the secret key used to validate incoming frames.
	// // Non signed frames are discarded, as well as frames with a version < 2.0.
	// InKey *frame.V2Key

	// Mavlink version used to encode messages. See Version
	// for the available options.
	// OutVersion Version
	// // the system id, added to every outgoing frame and used to identify this
	// // node in the network.
	// OutSystemID byte
	// // (optional) the component id, added to every outgoing frame, defaults to 1.
	// OutComponentID byte
	// (optional) the secret key used to sign outgoing frames.
	// This feature requires a version >= 2.0.
	// OutKey *frame.V2Key

	// (optional) automatically request streams to detected Ardupilot devices,
	// that need an explicit request in order to emit telemetry stream.
	// StreamRequestEnable bool
	// // (optional) the requested stream frequency in Hz. It defaults to 4.
	// StreamRequestFrequency int

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
