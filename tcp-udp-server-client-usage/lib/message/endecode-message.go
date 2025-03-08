package message

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"

	"github.com/sigurn/crc16"

	. "server-transport-go-usage/lib/utils"
)

const (
	PKG_START_FLAG = uint8(0xFE)
)

// protocal format is : 8 byte head + payload + 2 byte(payload crc16)
type HeaderMessage struct {
	StartFlag  uint8  // 消息起始位
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

	header[0] = m.HeaderMessage.StartFlag
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
	LogPrintf("encoded whole package len: %v\n", n)
	return n, nil
}

// UnPackageMessage 读取数据，解析协议，获取要每个有效的包。
func (m *DecodedMessage) UnPackageMessage(scanner *bufio.Scanner, rw *ReadWriter) error {
	if scanner.Scan() {
		err := m.ParseAndValidPkg(scanner.Bytes())
		if err != nil {
			LogPrintf("无效数据包: %v\n", err)
			return err
		}
		// 解包，返回一个 特性类型的 数据指针。
		if err := m.Unpackage(rw); err != nil {
			LogPrintln("parse message fail, err: ", err)
			return err
		}
		return nil
	}

	if err := scanner.Err(); err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			//if errors.Is(err, os.ErrDeadlineExceeded) {
			LogPrintln("it timeout.")
			return newIoTimeoutError("time out")
		}

		LogPrintln("scan request message fail: ", err)
		return err
	}
	return newEmptyPkgError("skip empty pkg")
}

func (m *DecodedMessage) Unpackage(rw *ReadWriter) error {
	// 解析： 使用 payload  和  msgType 进行解析
	msgData := rw.AllocateMsgData(m.PkgType)
	codeHandle := rw.GetCodecs()
	m.DecodedMsg = msgData
	err := codeHandle.Decode(m.payloadBin, msgData)
	if err != nil {
		return newError("decode payload bin fail, err: %v", err)
	}
	return nil
}

// 解析并校验数据包； 返回有效的包：
func (m *DecodedMessage) ParseAndValidPkg(pkt []byte) error {
	pkgLen := len(pkt)
	// 基础长度校验
	if pkgLen < minPacket {
		LogPrintf("receive pkg msg len: %d, less head len.", pkgLen)
		return newError("package received len: %v less than: %v; 数据包长度不足", pkgLen, minPacket)
	}

	// 分解数据包
	header := pkt[:headerSize]
	payload := pkt[headerSize : pkgLen-crcSize]
	storedCRC := binary.BigEndian.Uint16(pkt[pkgLen-crcSize:])

	// 校验payload长度
	payloadLen := binary.BigEndian.Uint16(header[1:3])
	if int(payloadLen) != len(payload) {
		LogPrintf("payload长度不匹配 (头声明:%d 实际:%d)\n", payloadLen, len(payload))
		return newError("payload data len: %v not eq payload len field: %v", len(payload), int(payloadLen))
	}

	// CRC校验（仅payload）
	computedCRC := crc16.Checksum(payload, crcTable)
	if computedCRC != storedCRC {
		LogPrintf("CRC校验失败 (预期:0x%04X 实际:0x%04X)\n", storedCRC, computedCRC)
		return newError("crc field: %x not eq calc value: %x", storedCRC, storedCRC)
	}

	// 解析其他头字段
	seq := binary.BigEndian.Uint16(header[3:5])
	source := header[5]
	msgType := binary.BigEndian.Uint16(header[6:8])

	// 输出结果
	LogPrintf("有效数据包: seq=%d source=0x%02X type=0x%04X payload_len=%d \n", seq, source, msgType, len(payload))

	m.HeaderMessage = &HeaderMessage{
		// 消息起始位
		StartFlag:  uint8(header[0]),                     //
		PayLoadLen: payloadLen,                           // 消息体长度 2
		PkgSeq:     binary.BigEndian.Uint16(header[3:5]), //每次发送消息的序列号
		DevType:    int8(header[5]),                      // 发送方类型
		PkgType:    binary.BigEndian.Uint16(header[6:8]), // 消息类型
	}

	m.payloadBin = payload
	return nil
}

func (m *DecodedMessage)parsePackage(reader *bufio.Reader)  (*HeaderMessage, []byte, error) {
	for {
		// 1. 寻找起始标志0xFD
		for {
			// 尝试读取一个字节，不推进Reader位置
			peekBytes, err := reader.Peek(1)
			if err != nil {
				if err == io.EOF {
					LogPrintf("peek fail, err: %v", err)
					return nil, nil, err // 数据不足，需重试
				}
				if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
					//if errors.Is(err, os.ErrDeadlineExceeded) {
					LogPrintln("it timeout.")
					return nil,nil, newIoTimeoutError("time out")
				}
				return nil, nil, fmt.Errorf("peak start flag, 读取错误: %v", err)
			}

			if peekBytes[0] == PKG_START_FLAG {
				break // 找到起始标志
			}

			// 丢弃非起始标志的字节
			_, _ = reader.Discard(1)
		}

		// 2. 检查是否有完整的8字节头部
		if _, err := reader.Peek(headerSize); err != nil {
			if err == io.EOF {
				LogPrintf("peek header fail, err: %v", err)
				return nil, nil, err // 数据不足，需等待
			}
			return nil, nil, fmt.Errorf("peak header, 读取头部失败: %v", err)
		}

		// 3. 读取完整头部
		headerBytes := make([]byte,headerSize)
		if _, err := io.ReadFull(reader, headerBytes); err != nil {
			return nil, nil, fmt.Errorf("读取头部失败: %v", err)
		}

		// 4. 解析头部（大端序）
		hdr := &HeaderMessage{
			StartFlag:  headerBytes[0],
			PayLoadLen: binary.BigEndian.Uint16(headerBytes[1:3]),
			PkgSeq:     binary.BigEndian.Uint16(headerBytes[3:5]),
			DevType:    int8(headerBytes[5]),
			PkgType:    binary.BigEndian.Uint16(headerBytes[6:8]),
		}

		// 5. 计算总需数据长度（头部+Payload+CRC）
		totalNeeded := int(hdr.PayLoadLen) + crcSize
		if _, err := reader.Peek(totalNeeded); err != nil {
			// 数据不足，需等待（此处已消耗头部，无法回退）
			return nil, nil, fmt.Errorf("peak payload and crc16, 等待Payload和CRC: %v", err)
		}

		// 6. 读取Payload和CRC
		payload := make([]byte, hdr.PayLoadLen)
		if _, err := io.ReadFull(reader, payload); err != nil {
			return nil, nil, fmt.Errorf("读取Payload失败: %v", err)
		}

		crcBytes := make([]byte, crcSize)
		if _, err := io.ReadFull(reader, crcBytes); err != nil {
			return nil, nil, fmt.Errorf("读取CRC失败: %v", err)
		}

		// 7. 校验CRC（Payload）
		receivedCRC := binary.BigEndian.Uint16(crcBytes)
		computedCRC := crc16.Checksum(payload, crc16.MakeTable(crc16.CRC16_CCITT_FALSE))
		if computedCRC != receivedCRC {
			// CRC失败，继续寻找下一个起始标志
			return nil, nil, errors.New("crc checksum is not right.")
		}

		// 8. 返回有效数据
		return hdr, payload, nil
	}
}
func( m*DecodedMessage) UnPackageMessageV2(reader *bufio.Reader, rw *ReadWriter) error {
	h, body, err := m.parsePackage(reader)
	if err != nil {
		return err
	}
	if h == nil {
		return fmt.Errorf("parsed fail.")
	}
	m.payloadBin = body
	m.HeaderMessage = h

	if h.PayLoadLen == 0 {
		m.DecodedMsg = nil
		m.payloadBin = nil
		return nil
	}

		// 解析： 使用 payload  和  msgType 进行解析
	msgData := rw.AllocateMsgData(m.PkgType)
	codeHandle := rw.GetCodecs()
	m.DecodedMsg = msgData
	err = codeHandle.Decode(m.payloadBin, msgData)
	if err != nil {
		return newError("decode payload bin fail, err: %v", err)
	}
	return nil
}
