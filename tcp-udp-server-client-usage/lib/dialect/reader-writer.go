package dialect

// import (
//
//
//
//
//
//
//

// 	"server-transport-go-usage/lib/codecs"
// 	"server-transport-go-usage/lib/message"
// )

// // ReadWriter is a Dialect Reader and Writer.
// type ReadWriter struct {
// 	// messageRWs map[uint32]*message.ReadWriter
// 	///
// 	messagePackage map[uint16] *message.CmdMsgReadWriter
// 	codeHandle codecs.BaseCodecs
// }
// func (rw *ReadWriter) GetCodecs() codecs.BaseCodecs {
// 	return rw.codeHandle
// }
// func (rw *ReadWriter) AllocateMsgData(cmd uint16) any {
// 	item, ok := rw.messagePackage[cmd]
// 	if !ok {
// 		return nil
// 	}

// 	return item.NewElemItem(cmd)
// }

// // NewReadWriter allocates a ReadWriter.
// func NewReadWriter(d *Dialect) (*ReadWriter, error) {
// 	rw := &ReadWriter{
// 		// messageRWs: make(map[uint32]*message.ReadWriter),
// 		messagePackage: make(map[uint16]*message.CmdMsgReadWriter),
// 		codeHandle:  new(codecs.ProtoCodecs),
// 	}
// 	for _, m := range d.CmdMsg {
// 		if m ==  nil {
// 			continue
// 		}
// 		if _, ok := rw.messagePackage[m.Cmd]; ok {
// 			continue
// 		}

// 		de, err := message.NewCmdMsgReadWriter(m.Cmd, m.MsgType)
// 		if err != nil {
// 			continue
// 		}
// 		rw.messagePackage[m.Cmd] = de
// 	}

// 	// for _, m := range d.Messages {
// 	// 	if _, ok := rw.messageRWs[m.GetID()]; ok {
// 	// 		return nil, fmt.Errorf("duplicate message with id %d", m.GetID())
// 	// 	}

// 	// 	de, err := message.NewReadWriter(m)
// 	// 	if err != nil {
// 	// 		return nil, fmt.Errorf("message %T: %w", m, err)
// 	// 	}

// 	// 	rw.messageRWs[m.GetID()] = de
// 	// }

// 	return rw, nil
// }

// func (rw *ReadWriter) GetCmdMessage(id uint16) *message.CmdMsgReadWriter {
// 	cmsg, ok := rw.messagePackage[id]
// 	if !ok {
// 		return nil
// 	}
// 	return cmsg
// }

// // // GetMessage returns the ReadWriter of a message.
// // func (rw *ReadWriter) GetMessage(id uint32) *message.ReadWriter {
// // 	mrw, ok := rw.messageRWs[id]
// // 	if !ok {
// // 		return nil
// // 	}
// // 	return mrw
// // }

