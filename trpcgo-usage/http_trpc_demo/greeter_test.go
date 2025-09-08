package main

import (
	"context"
	"reflect"
	"testing"

	pb "helloworld/pb"

	"go.uber.org/mock/gomock"
	_ "trpc.group/trpc-go/trpc-go/http"
)

//go:generate go mod tidy
//go:generate mockgen -destination=stub/helloworld/pb/helloworld_mock.go -package=pb -self_package=helloworld/pb --source=stub/helloworld/pb/helloworld.trpc.go

func Test_greeterImpl_Hello(t *testing.T) {
	// Start writing mock logic.
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	greeterService := pb.NewMockGreeterService(ctrl)
	var inorderClient []*gomock.Call
	// Expected behavior.
	m := greeterService.EXPECT().Hello(gomock.Any(), gomock.Any()).AnyTimes()
	m.DoAndReturn(func(ctx context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
		s := &greeterImpl{}
		return s.Hello(ctx, req)
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
			if rsp, err = greeterService.Hello(tt.args.ctx, tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("greeterImpl.Hello() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(rsp, tt.args.rsp) {
				t.Errorf("greeterImpl.Hello() rsp got = %v, want %v", rsp, tt.args.rsp)
			}
		})
	}
}
