package example

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/achilsh/go-dev-tools-collection/base-lib/config/file"
	logger "github.com/achilsh/go-dev-tools-collection/base-lib/log"
	logctx "github.com/achilsh/go-dev-tools-collection/base-lib/log/log_context"
)
type DB struct {
	Host string `yaml: "host"`
	Port int `yaml: "port"`
}
type Redis struct {
	Host string `yaml:"host"`
	Port int `yaml:"port"`
}
type DemoConfig struct {
	DB *DB  `yaml:"db"` // 要求字段名和类型 字母要一样，否则解析解析不出来
	Redis *Redis  `yaml:"redis"` // 要求字段名和类型 字母要一样，否则解析解析不出来
}


func Parsefile(t *testing.T) {
	fileName := "./config.yaml"
	cfg := file.NewConfig(fileName)
	var cfgItem DemoConfig
	err := cfg.Init(&cfgItem)
	if err != nil {
		fmt.Println("init config error: ", err)
		return 
	}

	fmt.Println("item: ", cfgItem.DB)
	assert.Equal(t, cfgItem.DB.Host, "0.0.0.0")
	assert.Equal(t, cfgItem.Redis.Port, 234)
}

func TestConfigParse(t *testing.T) {
	Parsefile(t)
}

func TestLogInit(t *testing.T) {
	logFIleName := "config.yaml"
	if err := logger.Init(logFIleName, "log"); err != nil {
		panic(err)
	}
	logger.Debug("this is debug log.")
	logger.Infof("this is info log")
	logger.Accessf("sdfadfad:====> %v",12)

	//
	{
		ctx := context.Background()
		logger.AccessCtx(ctx, "no with request id and uerId ")
		logger.AccessfCtx(ctx, "this is demo: %v", "no with  id and userid")
	}

	{
		ctx := context.Background()
		ctx = logctx.WithRequestID(ctx)
		logger.AccessCtx(ctx, "1111111111111111")
		logger.AccessfCtx(ctx, "this is demo: %v", 1231313)
	}

	{
		ctx := context.Background()
		ctx = logctx.WithRequestID(ctx)
		ctx = logctx.WithCtxUserID(ctx, "0123-i231nsld-13i1313")
		logger.AccessCtx(ctx, "333333333333333333333")
		logger.AccessfCtx(ctx, "this is demo....: %v", 565656565)
	}

	{
		ctx := context.Background()
		logger.InfoCtx(ctx, "this is info demo without user and request id.")
	}

	{
		ctx := context.Background()
		ctx = logctx.WithRequestID(ctx)
		logger.InfoCtx(ctx, "this is info demo without user and with request id.")
	}

	{
		ctx := context.Background()
		ctx = logctx.WithCtxUserID(ctx, "0000000000")
		logger.InfoCtx(ctx, "this is info demo  user and without request id.")
	}

	{
		ctx := context.Background()
		ctx = logctx.WithRequestID(ctx, "999999999999999999")
		ctx = logctx.WithCtxUserID(ctx, "0000000000")
		logger.InfoCtx(ctx, "this is info demo  user and with request id.")
	}
	{
		ctx := context.Background()
		logger.InfofCtx(ctx, "this is format with,  without id and without userid, value: %v", "----+++======")
	}
	{
		ctx := context.Background()
		ctx = logctx.WithRequestID(ctx)
		logger.InfofCtx(ctx, "this is format with,  id and without userid, value: %v", "kkkkkkkkkkkkkkkkkkk")
	}

	{
		ctx := context.Background()
		ctx = logctx.WithCtxUserID(ctx, "user_Id: 88888888888")
		logger.InfofCtx(ctx, "this is format with, without id and userid, value: %v", "jjjjjjjjjjjjjjjj")
	}

	{
		ctx := context.Background()
		ctx = logctx.WithCtxUserID(ctx, "user_Id: ---------------000==0-")
		ctx = logctx.WithRequestID(ctx)
		logger.InfofCtx(ctx, "this is format with, with id and userid, value: %v", "iiiiiiiiiiiiiiiiiiiiii")
	}
}
