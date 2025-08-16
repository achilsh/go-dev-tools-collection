package logic

import (
	"context"
	"fmt"

	"http_demo_server/internal/svc"
	"http_demo_server/internal/types"

	pb "rpc_demo_server1/rpc_demo_server1"
	mp2 "rpc_demo_server2/rpc_demo_server2"

	"github.com/zeromicro/go-zero/core/logx"
)

type Http_demo_serverLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewHttp_demo_serverLogic(ctx context.Context, svcCtx *svc.ServiceContext) *Http_demo_serverLogic {
	return &Http_demo_serverLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *Http_demo_serverLogic) Http_demo_server(req *types.Request) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line
	// call rpc clients.

	rpcClientReq := &pb.Request{
		Ping: fmt.Sprintf("rpc req msg: %v", req.Name),
	}

	rsp, err := l.svcCtx.DirectConnClient.Ping(context.Background(), rpcClientReq)
	if err != nil {
		l.Logger.Errorf("rpc client call fail, err: %v", err)
		return nil, err
	}

	if rsp == nil {
		l.Logger.Errorf("rpc client call response is nil ")
		return nil, fmt.Errorf("rpc client call response is nil")
	}
	resp = &types.Response{
		Message: fmt.Sprintf("rpc response data: %v", rsp.Pong),
	}

	r, e := l.call_more_rpc_services(req)
	if e != nil {
		return nil, e
	}

	resp.Message += ", " + r.Message
	return resp, nil
}

func (l *Http_demo_serverLogic) call_more_rpc_services(req *types.Request) (resp *types.Response, err error) {

	rpcClientReq := &mp2.Request{
		Ping: fmt.Sprintf("rpc req msg: %v", req.Name),
	}

	rsp, err := l.svcCtx.MoreCliConnClients.Ping(context.Background(), rpcClientReq)
	if err != nil {
		l.Logger.Errorf("rpc client call fail, err: %v", err)
		return nil, err
	}

	if rsp == nil {
		l.Logger.Errorf("rpc client call response is nil ")
		return nil, fmt.Errorf("rpc client call response is nil")
	}
	resp = &types.Response{
		Message: fmt.Sprintf("more rpc client: %v", rsp.Pong),
	}

	return resp, nil
}
