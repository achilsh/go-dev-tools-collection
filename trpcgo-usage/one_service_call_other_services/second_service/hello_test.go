package main

import (
	"context"
	"reflect"
	"testing"

	pb "second_demo/helloworld"

	"go.uber.org/mock/gomock"
	_ "trpc.group/trpc-go/trpc-go/http"
)

//go:generate go mod tidy
//go:generate mockgen -destination=stub/second_demo/helloworld/second_server_mock.go -package=helloworld -self_package=second_demo/helloworld --source=stub/second_demo/helloworld/second_server.trpc.go

func Test_helloImpl_SayHi(t *testing.T) {
	// Start writing mock logic.
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	helloService := pb.NewMockHelloService(ctrl)
	var inorderClient []*gomock.Call
	// Expected behavior.
	m := helloService.EXPECT().SayHi(gomock.Any(), gomock.Any()).AnyTimes()
	m.DoAndReturn(func(ctx context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
		s := &helloImpl{}
		return s.SayHi(ctx, req)
	})
	gomock.InOrder(inorderClient...)

	// Start writing unit test logic.
	type args struct {
		ctx context.Context
		req *pb.HelloRequest
		rsp *pb.HelloReply
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var rsp *pb.HelloReply
			var err error
			if rsp, err = helloService.SayHi(tt.args.ctx, tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("helloImpl.SayHi() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(rsp, tt.args.rsp) {
				t.Errorf("helloImpl.SayHi() rsp got = %v, want %v", rsp, tt.args.rsp)
			}
		})
	}
}
