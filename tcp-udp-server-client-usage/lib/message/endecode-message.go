package message

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"

	"github.com/sigurn/crc16"
	. "server-transport-go-usage/lib/utils"
)

const (
	PKG_START_FLAG = 0xFD
)

// protocal format is : 8 byte head + payload + 2 byte(payload crc16)
type HeaderMessage struct {
	StartFlag  int8   // 消息起始位
	PayLoadLen uint16 // 消息体长度
	PkgSeq     uint16 //每次发送消息的序列号
	DevType    int8   // 发送方类型
	PkgType    uint16 // 消息类型
}

// DecodedMessage 解析成具体业务的消息; need implement MsgFrame interface.
type DecodedMessage struct {
	*HeaderMessage
	DecodedMsg any    // payload 实际业务协议数据的结构体 变量指针
	payloadBin []byte // payload 二进制数据
}

func (m *DecodedMessage) GetPayLoad() uint16 {
	return m.HeaderMessage.PayLoadLen
}

func (m *DecodedMessage) GetPkgSeq() uint16 {
	return m.HeaderMessage.PkgSeq
}

func (m *DecodedMessage) GetDevType() int8 {
	return m.HeaderMessage.DevType
}
func (m *DecodedMessage) GetPkgType() uint16 {
	return m.HeaderMessage.PkgType
}

// 返回实际的消息体的结构体指针
func (m *DecodedMessage) GetMessage() any {
	return m.DecodedMsg
}
func (m *DecodedMessage) SetMessage(content []byte) {
	m.payloadBin = content
}

// SetPayLoadLen 设置二进制payload的长度
func (m *DecodedMessage) SetPayLoadLen(clen uint16) {
	m.PayLoadLen = clen
}

// 把payload 序列化的 数据 拼接成完整包，并序列化二进制；其中 header值在调用该函数前已经部分设置
func (m *DecodedMessage) PackageMessage(buf []byte) (int, error) {
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

// UnPackageMessage 读取数据，解析协议，获取要每个有效的包。
func (m *DecodedMessage) UnPackageMessage(scanner *bufio.Scanner, rw *ReadWriter) error {
	if scanner.Scan() {
		err := m.parseAndValidPkg(scanner.Bytes())
		if err != nil {
			LogPrintf("无效数据包: %v\n", err)
			return err
		}
		// 解包，返回一个 特性类型的 数据指针。
		if err := m.Unpackage(rw); err != nil {
			LogPrintln("parse message fail, err: ", err)
			return err
		}
	}

	if err := scanner.Err(); err != nil {
		LogPrintln("scan request message fail: ", err)
		return err
	}
	return nil
}

func (m *DecodedMessage) Unpackage(rw *ReadWriter) error {
	// 解析： 使用 payload  和  msgType 进行解析
	msgData := rw.AllocateMsgData(m.PkgType)
	codeHandle := rw.GetCodecs()
	return codeHandle.Decode(m.payloadBin, msgData)
}

// 解析并校验数据包； 返回有效的包：
func (m *DecodedMessage) parseAndValidPkg(pkt []byte) error {
	pkgLen := len(pkt)
	// 基础长度校验
	if pkgLen < minPacket {
		LogPrintf("receive pkg msg len: %d, less head len\n", pkgLen)
		return errors.New("数据包长度不足")
	}

	// 分解数据包
	header := pkt[:headerSize]
	payload := pkt[headerSize : pkgLen-crcSize]
	storedCRC := binary.BigEndian.Uint16(pkt[pkgLen-crcSize:])

	// 校验payload长度
	payloadLen := binary.BigEndian.Uint16(header[1:3])
	if int(payloadLen) != len(payload) {
		LogPrintf("payload长度不匹配 (头声明:%d 实际:%d)\n", payloadLen, len(payload))

		return errors.New("pkg payload len not eq len field.")
	}

	// CRC校验（仅payload）
	computedCRC := crc16.Checksum(payload, crcTable)
	if computedCRC != storedCRC {
		LogPrintf("CRC校验失败 (预期:0x%04X 实际:0x%04X)\n", storedCRC, computedCRC)
		return errors.New("received pkg crc checksum not eq.")
	}

	// 解析其他头字段
	seq := binary.BigEndian.Uint16(header[3:5])
	source := header[5]
	msgType := binary.BigEndian.Uint16(header[6:8])

	// 输出结果
	LogPrintf("有效数据包: seq=%d source=0x%02X type=0x%04X payload_len=%d\n",
		seq, source, msgType, len(payload))

	m.HeaderMessage = &HeaderMessage{
		// 消息起始位
		StartFlag:  int8(header[0]),                      //
		PayLoadLen: payloadLen,                           // 消息体长度 2
		PkgSeq:     binary.BigEndian.Uint16(header[3:5]), //每次发送消息的序列号
		DevType:    int8(header[5]),                      // 发送方类型
		PkgType:    binary.BigEndian.Uint16(header[6:8]), // 消息类型
	}

	m.payloadBin = payload
	return nil
}
