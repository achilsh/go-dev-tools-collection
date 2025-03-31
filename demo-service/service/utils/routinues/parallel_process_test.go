package routinues

import (
	"demo-service/service/utils/mock/mock_log"
	"testing"
	"time"

	logger "github.com/achilsh/go-dev-tools-collection/base-lib/log"
)

func TestParallelProcess(t *testing.T) {
	mock_log.LoggerMock()

	func() {
		beginTm := time.Now()
		defer func() {
			logger.Infof("cost tm: %v", time.Now().Sub(beginTm))
		}()
		WrapperFnWithTimeout(1000, func() {
			time.Sleep(2 * time.Second)
		})
	}()
	logger.Infof("wait other task...")

	time.Sleep(5 * time.Second)

}
