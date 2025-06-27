package main

import (
	"github.com/achilsh/go-dev-tools-collection/demo-service/service/utils/mock/mock_log"
	gobuildopt "github.com/achilsh/go-dev-tools-collection/go-build-opt"
)

func main() {
	mock_log.LoggerMock()
	gobuildopt.ShowVersion()
}
