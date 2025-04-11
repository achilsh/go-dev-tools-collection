package example

import (
	"testing"

	logger "github.com/achilsh/go-dev-tools-collection/base-lib/log"
	"github.com/achilsh/go-dev-tools-collection/demo-service/service/utils/mock/mock_log"
)

func TestCallOne(t *testing.T) {
	mock_log.LoggerMock()
	logger.Infof("ssss")
}
