package message

import (
	"encoding/binary"

	"github.com/sigurn/crc16"
)

const (
	headerSize = 8 // 1(magic) + 2(payloadLen) + 2(seq) + 1(source) + 2(type)
	crcSize    = 2
	minPacket  = headerSize + crcSize // 最小有效包长度（空payload情况）
)

var crcTable = crc16.MakeTable(crc16.CRC16_CCITT_FALSE)

// ProtocolSplitter 自定义协议分割函数（支持错误恢复）
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
		if data[i] != PKG_START_FLAG {
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
