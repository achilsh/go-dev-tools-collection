package tests

import (
	"context"
	"testing"

	_ "github.com/achilsh/go-dev-tools-collection/aws_usages/lambda_usage/second_demo/log_lib"
	loglib "github.com/achilsh/go-dev-tools-collection/aws_usages/lambda_usage/second_demo/log_lib"
	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/logx"
)

func TestBaseLog(t *testing.T) {
	defer loglib.DeconstructLog()

	for i := 0; i < 3; i++ {
		logx.Debug("This is a debug message")
		logx.Debugf("This is a formatted debug message: %s", "debug details")

		logx.Info("this is info message.")
		logx.Infof("this is formatted info message: %s", "info details")

		logx.Error("this is error message.")
		logx.Errorf("this is formatted error message: %s", "error details")

		logx.ErrorStack("this is error stack message.")
		logx.ErrorStackf("this is formatted error stack message: %s", "error stack details")

	}

	for i := 0; i < 3; i++ {
		ctx := context.Background()
		ctx = context.WithValue(ctx, "traceid", "1234567890")
		logc.Debug(ctx, "this is demo debug message.")

	}
}
