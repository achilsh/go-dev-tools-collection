package frame

import (
	"bufio"

	"server-transport-go-usage/lib/message"
)

type MsgFrame interface {
	GetPayLoad() uint16
	GetPkgSeq() uint16
	GetDevType() int8
	GetPkgType() uint16
	// 返回实际的消息体的结构体指针
	GetMessage() any
	//
	SetMessage(content []byte)
	SetPayLoadLen(clen uint16)
	// 把消息体和头部 打包成二进制数据
	PackageMessage(buf []byte) (int, error)
	// 读数据，解具体协议
	UnPackageMessage(*bufio.Scanner, *message.ReadWriter) error
}
