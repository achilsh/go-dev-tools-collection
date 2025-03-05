package lib

import (
	"server-transport-go-usage/lib/frame"
)

// Event is the interface implemented by all events received with node.Events().
type Event interface {
	isEventOut()
}

// EventChannelOpen is the event fired when a channel gets opened.
type EventChannelOpen struct {
	Channel *Channel
}

func (*EventChannelOpen) isEventOut() {}

// EventChannelClose is the event fired when a channel gets closed.
type EventChannelClose struct {
	Channel *Channel
}

func (*EventChannelClose) isEventOut() {}

// EventFrame is the event fired when a frame is received.
type EventFrame struct {
	// the frame
	// Frame frame.Frame
	Frame frame.MsgFrame

	// the channel from which the frame was received
	Channel *Channel
}

func (*EventFrame) isEventOut() {}

// EventParseError is the event fired when a parse error occurs.
type EventParseError struct {
	// the error
	Error error

	// the channel used to send the frame
	Channel *Channel
}

func (*EventParseError) isEventOut() {}

// EventStreamRequested is the event fired when an automatic stream request is sent.
type EventStreamRequested struct {
	// the channel to which the stream request is addressed
	Channel *Channel
	// the system id to which the stream requests is addressed
	SystemID byte
	// the component id to which the stream requests is addressed
	ComponentID byte
}

func (*EventStreamRequested) isEventOut() {}
