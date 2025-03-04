package message

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"errors"

	"github.com/sigurn/crc16"

	"server-transport-go-usage/lib/dialect"
)

const (
	headerSize = 8 // 1(magic) + 2(payloadLen) + 2(seq) + 1(source) + 2(type)
	magicByte  = 0xFE
	crcSize    = 2
	minPacket  = headerSize + crcSize // 最小有效包长度（空payload情况）
)

var crcTable = crc16.MakeTable(crc16.CRC16_CCITT_FALSE)

// 自定义协议分割函数（支持错误恢复）
func ProtocolSplitter(data []byte, atEOF bool) (advance int, token []byte, err error) {
	dataLen := len(data)

	// 遍历数据寻找合法起始位置
	for i := 0; i <= dataLen; {
		// 剩余数据不足最小包头
		if dataLen-i < headerSize {
			if atEOF {
				// 丢弃无法解析的剩余数据
				return dataLen, nil, nil
			}
			// 保留当前数据等待后续输入
			return i, nil, nil
		}

		// 检查 Magic 字节有效性
		if data[i] != magicByte {
			i++ // 尝试下一个字节
			continue
		}

		// 解析 payload 长度
		payloadLen := int(binary.BigEndian.Uint16(data[i+1 : i+3]))
		totalSize := headerSize + payloadLen + crcSize

		// 完整性检查
		if dataLen-i >= totalSize {
			// 返回完整数据包
			return i + totalSize, data[i : i+totalSize], nil
		}

		// 数据不完整但后续可能还有数据
		if !atEOF {
			return i, nil, nil // 等待更多数据
		}

		// 已到末尾但数据不完整
		return dataLen, nil, nil // 丢弃不完整数据
	}

	return dataLen, nil, nil // 遍历完未找到有效包
}


// 解析并校验数据包； 返回有效的包：
func parseAndValidate(pkt []byte) (error, *UnDecodedMessage) {
        // 基础长度校验
        if len(pkt) < minPacket {
                return errors.New("数据包长度不足"),nil
        }

        // 分解数据包
        header := pkt[:headerSize]
        payload := pkt[headerSize : len(pkt)-crcSize]
        storedCRC := binary.BigEndian.Uint16(pkt[len(pkt)-crcSize:])

        // 校验payload长度
        payloadLen := binary.BigEndian.Uint16(header[1:3])
        if int(payloadLen) != len(payload) {
            return fmt.Errorf("payload长度不匹配（头声明:%d 实际:%d）", payloadLen, len(payload)), nil
        }

        // CRC校验（仅payload）
        computedCRC := crc16.Checksum(payload, crcTable)
        if computedCRC != storedCRC {
                return fmt.Errorf("CRC校验失败（预期:0x%04X 实际:0x%04X）", storedCRC, computedCRC),nil
        }

        // 解析其他头字段
        seq := binary.BigEndian.Uint16(header[3:5])
        source := header[5]
        msgType := binary.BigEndian.Uint16(header[6:8])

        // 输出结果
        fmt.Printf("有效数据包: seq=%d source=0x%02X type=0x%04X payload_len=%d\n",
                seq, source, msgType, len(payload))

        retMessage := &HeaderMessage  {
			StartFlag: header[0:1]   // 消息起始位 1
			PayLoadLen : payloadLen  // 消息体长度 2
			PkgSeq: binary.BigEndian.Uint16(header[3:5]) //每次发送消息的序列号
			DevType:  header[5:6]// 发送方类型
			PkgType: binary.BigEndian.Uint16(header[6:8])// 消息类型
        }
        //
        var ret = &UnDecodedMessage {
            HeaderMessage:retMessage,
            PayLoad: payload,
        }
        return nil, ret
}

func Unpackage(msg *UnDecodedMessage, rw *ReadWriter) any {
	// add item.
	// 解析： 使用 payload  和  msgType 进行解析
	msgData := rw.AllocateMsgData(msg.PkgType)
	codeHandle := rw.GetCodecs()
	codeHandle(msg.PayLoad, msgData)

	// msgType 和 msgData 作为一个整体发给业务处理线程。
	return msgData

}

// ReadAndParseMessage 通过 io.reader 读取数据
func ReadAndParseMessage(scanner *bufio.Scanner, rw *ReadWriter) (*DecodedMessage, error) {
	var parsedMessage *DecodedMessage
	if scanner.Scan() {
		pkt := scanner.Bytes()
		err, pkgUncoded:= parseAndValidate(pkt);
		if err != nil {
			fmt.Printf( "无效数据包: %v\n", err)
			return nil, err
		}

		// 解包，返回一个 特性类型的 数据指针。
		codedMsg :=  Unpackage(pkgUncoded, rw)
		// 把解析好的包发送到业务处理协程中。
		decodedMsg := &DecodedMessage {
			HeaderMessage: pkgUncoded.HeaderMessage,
			DecodedMsg: codedMsg,
		}
		// return decodedMsg
		parsedMessage =  decodedMsg
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("扫描错误:", err)
		return nil, err
	}
	return parsedMessage, nil
}
