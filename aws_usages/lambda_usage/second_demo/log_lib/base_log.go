package loglib

import "github.com/zeromicro/go-zero/core/logx"

type LogConfType struct {
	Mode        string
	Level       string
	Encoding    string
	ServiceName string
}

var LogConfItem = LogConfType{
	Mode:        "console", // 输出模式，支持 console、file 和 volume 三种模式
	Level:       "debug",   // debug,info,error,severe
	Encoding:    "plain",   // json,plain
	ServiceName: "lambda-server",
}

func init() {
	logx.MustSetup(logx.LogConf{
		Mode:        LogConfItem.Mode,     // 输出模式，支持 console、file 和 volume 三种模式
		Path:        "logs",               // 文件路径，console 模式可忽略
		Level:       LogConfItem.Level,    // debug,info,error,severe
		KeepDays:    7,                    //only on Mode is file
		Encoding:    LogConfItem.Encoding, // json,plain
		Compress:    false,
		Stat:        true,
		ServiceName: LogConfItem.ServiceName,
	})
}

func DeconstructLog() {
	logx.Close()
}
