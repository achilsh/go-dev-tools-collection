package router_wrapper

import (
	"encoding/json"
	"testing"

	logger "github.com/achilsh/go-dev-tools-collection/base-lib/log"

	"demo-service/service/utils/mock/mock_log"
)

type DemoJson struct {
	Data string `json:"data"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func TestJsonToMap(t *testing.T) {
	mock_log.LoggerMock()
	var a = &DemoJson{
		Data: "hello",
		Name: "sz",
		Age:  12,
	}
	buf, _ := json.Marshal(a)

	var ret map[string]interface{}
	json.Unmarshal(buf, &ret)
	for k, v := range ret {
		logger.Infof("k: %v, v: %v\n", k, v)
	}

}
