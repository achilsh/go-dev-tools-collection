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

func Test_greeterImpl_SayHello(t *testing.T) {
	// Start writing mock logic.
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	greeterService := pb.NewMockGreeterService(ctrl)
	var inorderClient []*gomock.Call
	// Expected behavior.
	m := greeterService.EXPECT().SayHello(gomock.Any(), gomock.Any()).AnyTimes()
	m.DoAndReturn(func(ctx context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
		s := &greeterImpl{}
		return s.SayHello(ctx, req)
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
			if rsp, err = greeterService.SayHello(tt.args.ctx, tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("greeterImpl.SayHello() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(rsp, tt.args.rsp) {
				t.Errorf("greeterImpl.SayHello() rsp got = %v, want %v", rsp, tt.args.rsp)
			}
		})
	}
}
