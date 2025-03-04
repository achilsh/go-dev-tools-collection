package message

import (
	"bytes"
	"encoding/binary"

	"github.com/sigurn/crc16"
)

const (
	PKG_START_FLAG = 0xFD
)
type HeaderMessage struct {
	StartFlag  int8   // 消息起始位
	PayLoadLen uint16 // 消息体长度
	PkgSeq     uint16 //每次发送消息的序列号
	DevType    int8   // 发送方类型
	PkgType    uint16 // 消息类型
}

// UnDecodedMessage 未解析成具体业务的消息
type UnDecodedMessage struct {
	*HeaderMessage
	PayLoad []byte
}


// DecodedMessage 解析成具体业务的消息; need implement MsgFrame interface.
type DecodedMessage struct {
	*HeaderMessage
	DecodedMsg any // 是结构体化对象指针变量
	//
	payloadBin []byte
}

func (m *DecodedMessage) GetPayLoad() uint16 {
	return m.HeaderMessage.PayLoadLen
}

func (m *DecodedMessage)GetPkgSeq() uint16 {
	return m.HeaderMessage.PkgSeq
}

func (m *DecodedMessage)GetDevType() int8 {
	return m.HeaderMessage.DevType
}
func (m *DecodedMessage)GetPkgType() uint16 {
	return m.HeaderMessage.PkgType
}

//返回实际的消息体的结构体指针
func (m *DecodedMessage)GetMessage() any {
	return m.DecodedMsg
}
func (m *DecodedMessage) SetMessage(content []byte) {
	m.payloadBin = content
}
// SetPayLoadLen 设置二进制payload的长度
func (m *DecodedMessage)SetPayLoadLen(clen uint16) {
	m.PayLoadLen = clen
}
// 把payload 序列化的 数据 拼接成完整包，并序列化二进制；其中 header值在调用该函数前已经部分设置
func (m *DecodedMessage) PackageMessage(buf []byte) (int, error)  {
	header := make([]byte, headerSize)

	header[0] = byte(m.HeaderMessage.StartFlag)
	binary.BigEndian.PutUint16(header[1:3], m.HeaderMessage.PayLoadLen)
	binary.BigEndian.PutUint16(header[3:5], m.HeaderMessage.PkgSeq)
	header[5] = byte(m.HeaderMessage.DevType)
	binary.BigEndian.PutUint16(header[6:8], m.HeaderMessage.PkgType)

	//
	crc := crc16.Checksum(m.payloadBin, crcTable)
	// 组合完整数据包
	bufTmp := new(bytes.Buffer)
	bufTmp.Write(header)
	bufTmp.Write(m.payloadBin)
	binary.Write(bufTmp, binary.BigEndian, crc)
	n := copy(buf, bufTmp.Bytes())
	return n, nil
}
