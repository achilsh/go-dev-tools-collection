package message

import "reflect"

// CmdMsgReadWriter 消息类型和消息结构体类型
type CmdMsgReadWriter struct {
	elemType reflect.Type
	cmd uint16
}

// NewCmdMsgReadWriter 注册消息类型
func NewCmdMsgReadWriter(cmd uint16, elemType reflect.Type)(*CmdMsgReadWriter, error) {
	return &CmdMsgReadWriter{
		elemType: elemType,
		cmd: cmd,
	}, nil
}


// NewElemItem 根据消息类型和结构体类型创建消息结构体对象
func (crw *CmdMsgReadWriter) NewElemItem(cmd uint16) any {
	item := reflect.New(crw.elemType)
	return item.Interface()
}