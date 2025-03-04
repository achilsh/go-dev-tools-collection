package message

import (
	"reflect"

	"server-transport-go-usage/lib/codecs"
)

type CmdMessageInfo struct {
	Cmd     uint16
	MsgType reflect.Type
}




// ReadWriter is a Dialect Reader and Writer.
type ReadWriter struct {
	codeHandle codecs.BaseCodecs //编解码处理器
	messagePackage map[uint16] *CmdMsgReadWriter //消息类型集合
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
func NewReadWriter(d []*CmdMessageInfo) (*ReadWriter, error) {
	rw := &ReadWriter{	
		codeHandle : new(codecs.ProtoCodecs),
		messagePackage: make(map[uint16]*CmdMsgReadWriter),
	}
	for _, m := range d {
		if m ==  nil {
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

	// for _, m := range d.Messages {
	// 	if _, ok := rw.messageRWs[m.GetID()]; ok {
	// 		return nil, fmt.Errorf("duplicate message with id %d", m.GetID())
	// 	}

	// 	de, err := message.NewReadWriter(m)
	// 	if err != nil {
	// 		return nil, fmt.Errorf("message %T: %w", m, err)
	// 	}

	// 	rw.messageRWs[m.GetID()] = de
	// }

	return rw, nil
}

func (rw *ReadWriter) GetCmdMessage(id uint16) *CmdMsgReadWriter {
	cmsg, ok := rw.messagePackage[id]
	if !ok {
		return nil
	}
	return cmsg
}