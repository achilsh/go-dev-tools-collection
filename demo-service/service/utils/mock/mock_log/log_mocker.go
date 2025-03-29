package mock_log

import (
	"fmt"

	logger "github.com/achilsh/go-dev-tools-collection/base-lib/log"
)

var mockErrorf = func(template string, args ...interface{}) {
	fmt.Printf(template, args...)
	return
}
var mockError = func(v ...interface{}) {
	fmt.Println(v)
	return
}
var mockErrorw = func(msg string, keysAndValues ...interface{}) {
	fmt.Print(msg)
	fmt.Println(keysAndValues...)
	return
}

func LoggerMock() {
	logger.Infof, logger.Debugf, logger.Errorf, logger.Warnf, logger.Panicf, logger.Fatalf = mockErrorf, mockErrorf, mockErrorf, mockErrorf, mockErrorf, mockErrorf
	logger.Info, logger.Debug, logger.Error, logger.Warn, logger.Panic, logger.Fatal = mockError, mockError, mockError, mockError, mockError, mockError
	logger.Infow, logger.Debugw = mockErrorw, mockErrorw
}
