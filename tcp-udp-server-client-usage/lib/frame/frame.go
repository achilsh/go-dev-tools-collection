package frame

import (
	"bufio"

	"server-transport-go-usage/lib/message"
)

// Frame is the interface implemented by frames of every supported version.
type Frame interface {
	// returns the system id of the author of the frame.
	GetSystemID() byte

	// returns the component id of the author of the frame.
	GetComponentID() byte

	// returns the sequence number in the frame
	GetSequenceNumber() byte

	// returns the message wrapped in the frame.
	GetMessage() message.Message

	// returns the checksum of the frame.
	GetChecksum() uint16

	// generates the checksum of the frame.
	GenerateChecksum(byte) uint16

	decode(*bufio.Reader) error
	encodeTo([]byte, []byte) (int, error)
}

type MsgFrame interface {
	GetPayLoad() uint16
	GetPkgSeq() uint16 
	GetDevType() int8 
	GetPkgType() uint16 
	//返回实际的消息体的结构体指针
	GetMessage() any
	//
	SetMessage(content []byte)
	SetPayLoadLen(clen uint16)

	PackageMessage(buf []byte) (int, error)
}
