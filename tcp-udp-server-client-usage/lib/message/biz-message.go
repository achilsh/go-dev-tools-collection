package message

import (
	"reflect"

	"server-transport-go-usage/lib/codecs"
)

type CmdMessageInfo struct {
	Cmd     uint16
	MsgType reflect.Type //is Pointer.
}

// ReadWriter 一个连接的 读写器
type ReadWriter struct {
	//连接上的编码器
	codeHandle codecs.BaseCodecs
	// 所有消息类型枚举和消息类型结构体抽象
	messagePackage map[uint16]*CmdMsgReadWriter
}

func (rw *ReadWriter) GetCodecs() codecs.BaseCodecs {
	return rw.codeHandle
}
func (rw *ReadWriter) AllocateMsgData(cmd uint16) any {
	item, ok := rw.messagePackage[cmd]
	if !ok {
		return nil
	}
	return item.NewElemItem(cmd)
}

// NewReadWriter allocates a ReadWriter.
func NewReadWriter(d map[uint16]*CmdMessageInfo) (*ReadWriter, error) {
	rw := &ReadWriter{
		codeHandle:     new(codecs.ProtoCodecs),
		messagePackage: make(map[uint16]*CmdMsgReadWriter),
	}
	for _, m := range d {
		if m == nil {
			continue
		}
		if _, ok := rw.messagePackage[m.Cmd]; ok {
			continue
		}

		de, err := NewCmdMsgReadWriter(m.Cmd, m.MsgType)
		if err != nil {
			continue
		}
		rw.messagePackage[m.Cmd] = de
	}

	return rw, nil
}

// GetCmdMessage 根据消息类型枚举获取消息结构体
func (rw *ReadWriter) GetCmdMessage(id uint16) *CmdMsgReadWriter {
	cmsg, ok := rw.messagePackage[id]
	if !ok {
		return nil
	}
	return cmsg
}
