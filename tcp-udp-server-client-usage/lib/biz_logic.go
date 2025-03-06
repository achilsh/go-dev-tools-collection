package lib

import (
	"context"

	. "server-transport-go-usage/lib/utils"
)

//
type FailLogicProcesser struct {

}


func (f *FailLogicProcesser) HandleParseFail(ctx context.Context, in Event) (error, any) {
	return nil,nil
}
func (f *FailLogicProcesser) HandleNewConn(ctx context.Context, in Event) (error, any) {
	LogPrintf("receive new connect, log new connect-------")
	return nil, nil
}

type MsgOneProcess struct {
	Cmd int

}
func NewMsgOneProcess() *MsgOneProcess{
	return &MsgOneProcess{
		Cmd:1,
	}
}
func (m* MsgOneProcess) Handle(ctx context.Context, in Event) (error, any) {
	LogPrintf("receive biz logic data, now it's processing.....")
	return nil, nil
}
//
