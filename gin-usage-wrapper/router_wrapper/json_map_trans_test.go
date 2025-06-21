package router_wrapper

import (
	"encoding/json"
	"testing"

	logger "github.com/achilsh/go-dev-tools-collection/base-lib/log"

	"github.com/achilsh/go-dev-tools-collection/gin-usage-wrapper/error_def"
	"github.com/achilsh/go-dev-tools-collection/gin-usage-wrapper/mock/mock_log"
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

type DemoResultData struct {
	A int    `json:"a"`
	B string `json:"b"`
}

func TestDumpResponse(t *testing.T) {
	mock_log.LoggerMock()
	responseData := &error_def.HttpResponse{
		ErrorMessage: "success",
		ErrorCode:    "",
	}
	respJson, _ := json.MarshalIndent(responseData, "", "  ")
	// logger.Infof("response: %+v", spew.Sdump(*responseData))

	logger.Infof("response: %+v", string(respJson))
}
