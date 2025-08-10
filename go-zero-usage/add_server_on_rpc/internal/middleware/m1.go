package middleware

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
)

func ExampleUnaryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp interface{}, err error) {
	// TODO: fill your logic here

	beginTim := time.Now()

	ret, err := handler(ctx, req)

	fmt.Println("cost time: ", time.Since(beginTim))
	return ret, err
}
