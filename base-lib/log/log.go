package log

import (
	logconf "ai-customer-service/service/utils/base_lib/log/config"
	"ai-customer-service/service/utils/base_lib/log/lumberjack"
	"ai-customer-service/service/utils/base_lib/log/zap"
	"ai-customer-service/service/utils/base_lib/log/zap/zapcore"
	"fmt"
	"github.com/go-yaml/yaml"
	"github.com/jinzhu/copier"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// / if in config file set level is "debug", add output to console for debug mode
// / config log file root node key is "logger"
const (
	config_level_debug_value = "debug"
	config_logger_root_key   = "logger"
)

// var logger *zap.SugaredLogger
type noWithParameterFunc func(v ...interface{})
type withParameterFunc func(template string, args ...interface{})

var Info, Debug, Error, Warn, Panic, Fatal func(v ...interface{})
var Infof, Debugf, Errorf, Warnf, Panicf, Fatalf func(template string, args ...interface{})
var Infow, Debugw func(msg string, keysAndValues ...interface{})

var logMap = map[zapcore.Level]noWithParameterFunc{}
var logMapF = map[zapcore.Level]withParameterFunc{}

// Access
func Access(v ...interface{}) {
	if fn, ok := logMap[zapcore.AccessLevel]; ok {
		if fn != nil {
			fn(v...)
		}
	}
}
func Accessf(template string, args ...interface{}) {
	if fn, ok := logMapF[zapcore.AccessLevel]; ok {
		if fn != nil {
			fn(template, args...)
		}
	}
}

// for specific log save directory by function callback
var updateLogSaveDir = ""

// for use in pad env, frontbackend read runtime env, and get log save directory
// then as Init function logSaveDir parament
func Init(logConfigFile string, logSaveDir ...string) error {
	// read config file
	yamlContent, err := ioutil.ReadFile(logConfigFile)
	if err != nil {
		return fmt.Errorf("read file %s failed, occured error: %s", logConfigFile, err)
	}

	if len(logSaveDir) > 0 {
		updateLogSaveDir = logSaveDir[0]
	}

	cfgObj := map[string]interface{}{}
	err = yaml.UnmarshalStrict(yamlContent, &cfgObj)
	if err != nil {
		return fmt.Errorf("Unmarshal file %s failed, occured error: %s", logConfigFile, err)
	}
	logConf := convLogConf(cfgObj[config_logger_root_key])

	// create log core
	var core zapcore.Core
	var accessCore zapcore.Core = nil
	core = zapcore.NewCore(
		getEncoder(logConf),
		getLogWriter(logConf.Rotate),
		//getLogWriter1(logConf.Rotate),
		getLevel(logConf),
	)

	if logConf.AccessRotate != nil {
		accessCore = zapcore.NewCore(
			getEncoder(logConf),
			getLogWriter(*(logConf.AccessRotate)),
			getLevel(logConf),
		)
	}

	if logConf.Level == config_level_debug_value {
		core = zapcore.NewTee(
			//for debug model: add output to console
			core, zapcore.NewCore(
				getEncoder(logConf), zapcore.Lock(os.Stdout), zapcore.DebugLevel),
		)
	}

	var logger *zap.SugaredLogger
	if len(strings.TrimSpace(logConf.FilterTopic)) > 0 {
		zapcore.SetHookedTopics(strings.Split(logConf.FilterTopic, ","))
		logger = zap.New(core, zap.AddCaller(), getHookOption()).Sugar()
	} else {
		logger = zap.New(core, zap.AddCaller()).Sugar()
	}

	var accessLogger *zap.SugaredLogger
	if accessCore != nil {
		accessLogger = zap.New(accessCore, zap.AddCaller(), zap.AddCallerSkip(2), zap.WithCaller(true)).Sugar() //zap.AddCallerSkip(2), zap.WithCaller(true)
		logMap[zapcore.AccessLevel] = accessLogger.Access
		logMapF[zapcore.AccessLevel] = accessLogger.Accessf
	}

	initGlobalFuncs(logger)
	return nil
}

func convStr(v interface{}, defaultVals ...string) string {
	if v == nil {
		if len(defaultVals) > 0 {
			return defaultVals[0]
		} else {
			return ""
		}

	} else {
		return v.(string)
	}
}

func convInt(v interface{}, defaultVals ...int) int {
	if v == nil {
		if len(defaultVals) > 0 {
			return defaultVals[0]
		} else {
			return 1
		}
	} else {
		return v.(int)
	}
}

func convBool(v interface{}, defaultVals ...bool) bool {
	if v == nil {
		if len(defaultVals) > 0 {
			return defaultVals[0]
		} else {
			return true
		}
	} else {
		return v.(bool)
	}
}

func NewRotateStruct(item logconf.RotateStruct) *logconf.RotateStruct {
	ret := logconf.RotateStruct{}
	copier.Copy(&ret, &item)
	return &ret
}
func convRotate(r interface{}) logconf.RotateStruct {
	rst := logconf.RotateStruct{}
	if r != nil {
		rm := r.(map[interface{}]interface{})
		for k, v := range rm {
			switch k.(string) {
			case "filename":
				rst.Filename = convStr(v, "app.log")
			case "maxsize":
				rst.Maxsize = convInt(v, 1)
			case "maxage":
				rst.Maxage = convInt(v, 7)
			case "maxbackups":
				rst.Maxbackups = convInt(v, 20)
			case "localtime":
				rst.Localtime = convBool(v)
			case "compress":
				rst.Compress = convBool(v)
			}
		}
	}
	if len(updateLogSaveDir) > 0 {
		fileName := filepath.Base(rst.Filename)
		rst.Filename = filepath.Join(updateLogSaveDir, fileName)
	}
	return rst
}

// read log config file to interface{} obj, and convert into logconf.LogConf object
func convLogConf(c interface{}) logconf.LogConf {
	rst := logconf.LogConf{}
	if c != nil {
		m := c.(map[interface{}]interface{})
		for k, v := range m {
			switch k.(string) {
			case "encoding":
				rst.Encoding = convStr(v, "text")
			case "level":
				rst.Level = convStr(v, "error")
			case "filterTopic":
				rst.FilterTopic = convStr(v, "")
			case "rotate":
				rst.Rotate = convRotate(v)
			case "accessRotate":
				rst.AccessRotate = NewRotateStruct(convRotate(v))
			}
		}
	}
	return rst
}

func initGlobalFuncs(l *zap.SugaredLogger) {
	Info = l.Info
	Infof = l.Infof
	Infow = l.Infow
	Debug = l.Debug
	Debugf = l.Debugf
	Debugw = l.Debugw
	Error = l.Error
	Errorf = l.Errorf
	Warn = l.Warn
	Warnf = l.Warnf
	Panic = l.Panic
	Panicf = l.Panicf
	Fatal = l.Fatal
	Fatalf = l.Fatalf
}

func getLogWriter(rotate logconf.RotateStruct) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   rotate.Filename,
		MaxSize:    rotate.Maxsize, //* 1024 * 1024, // 10MB
		MaxBackups: rotate.Maxbackups,
		MaxAge:     rotate.Maxage,
		LocalTime:  rotate.Localtime,
		Compress:   rotate.Compress,
	}
	return zapcore.AddSync(lumberJackLogger)
}

// for debug
func getLogWriter1(rotate logconf.RotateStruct) zapcore.WriteSyncer {
	file, _ := os.OpenFile(rotate.Filename, os.O_APPEND|os.O_CREATE, os.ModeAppend)
	return zapcore.AddSync(file)
}

func getEncoder(conf logconf.LogConf) zapcore.Encoder {
	encoderConfig := zapcore.EncoderConfig{
		MessageKey:    "msg",
		LevelKey:      "level",
		TimeKey:       "time",
		NameKey:       "name",
		CallerKey:     "caller",
		FunctionKey:   "func",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.LowercaseLevelEncoder,
		//EncodeTime:     zapcore.RFC3339TimeEncoder,
		EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02T15:04:05.000"),
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	switch strings.TrimSpace(conf.Encoding) {
	case "json":
		return zapcore.NewJSONEncoder(encoderConfig)
	case "text":
		return zapcore.NewConsoleEncoder(encoderConfig)
	default:
		return zapcore.NewConsoleEncoder(encoderConfig)
	}
}

var levelMap = map[string]zapcore.Level{
	"debug":  zapcore.DebugLevel,
	"info":   zapcore.InfoLevel,
	"warn":   zapcore.WarnLevel,
	"error":  zapcore.ErrorLevel,
	"panic":  zapcore.PanicLevel,
	"fatal":  zapcore.FatalLevel,
	"none":   zapcore.ErrorLevel,
	"access": zapcore.AccessLevel,
}

func getLevel(conf logconf.LogConf) zapcore.Level {
	if level, ok := levelMap[conf.Level]; ok {
		return level
	}
	return zapcore.ErrorLevel
}

func getHookOption() zap.Option {
	return zap.Hooks(
		func(e zapcore.Entry) error {
			return nil
		})
}
