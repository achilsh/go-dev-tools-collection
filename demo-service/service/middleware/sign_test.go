package middleware

import (
	"encoding/json"
	"fmt"
	"testing"

	logger "github.com/achilsh/go-dev-tools-collection/base-lib/log"
	"github.com/stretchr/testify/assert"

	"demo-service/service/utils/mock/mock_log"
)

type DemoJson struct {
	NonceStr        string `json:"nonce_str"`
	TimeStampSecond string `json:"timestamp"`
	Sign            string `json:"sign"`
	Uid             string `json:"uid"`
	AccessToken     string `json:"access_token"`
	CurrentPage     int    `json:"current_page"`
}

func TestBuildSign(t *testing.T) {
	mock_log.LoggerMock()

	var a = &DemoJson{
		NonceStr:        "QrZHJAlEf3mgG",
		TimeStampSecond: "1722846693",
		Sign:            "1adfadf",
		Uid:             "ce789bdf-1afb-4d7b-873a-7a886ffee870",
		AccessToken:     "sdfadfafsdfadsf",
		CurrentPage:     1,
	}
	buf, _ := json.Marshal(a)

	var ret map[string]any
	json.Unmarshal(buf, &ret)
	for k, v := range ret {
		fmt.Println(k, v)
	}

	signCalcData, err := calcSign(ret, "1113213")
	if err != nil {
		logger.Errorf("calc sign fail: %v", err)
	}
	logger.Infof("signe calc data: %v\n", signCalcData)

	errCode := CheckSign(ret, "1113213")
	if errCode != "" {
		logger.Errorf("check sign fail, err: %v\n", errCode)
	}
}

type Demo1Json struct {
	NonceStr        string `json:"nonce_str"`
	TimeStampSecond string `json:"timestamp"`
	Sign            string `json:"sign"`
	Uid             string `json:"uid"`
	AccessToken     string `json:"access_token"`
	ABC             int    `json:"abc"`
}
type Demo2Json struct {
	ABC int `json:"abc"`
}

func TestParser(t *testing.T) {
	var d1 = &Demo1Json{
		NonceStr:        "QrZHJAlEf3mgG",
		TimeStampSecond: "1722846693",
		Sign:            "1adfadf",
		Uid:             "sdfasdfa",
		AccessToken:     "sdfadfafsdf",
		ABC:             12345,
	}
	data, err := json.Marshal(d1)
	assert.Equal(t, err, nil)
	//使用简单简化数据来解析：
	var dst2 = &Demo2Json{}
	err = json.Unmarshal(data, dst2)
	assert.Equal(t, err, nil)
	assert.Equal(t, dst2.ABC, d1.ABC)
}

type Demo3Json struct {
	NonceStr        string `json:"nonce_str"`
	TimeStampSecond string `json:"timestamp"`
	Sign            string `json:"sign"`
	Uid             string `json:"uid"`
	AccessToken     string `json:"access_token"`
}
type Demo4Json struct {
}

func TestJsonParse2(t *testing.T) {
	var d1 = &Demo3Json{
		NonceStr:        "QrZHJAlEf3mgG",
		TimeStampSecond: "1722846693",
		Sign:            "1adfadf",
		Uid:             "sdfasdfa",
		AccessToken:     "sdfadfafsdf",
	}
	data, err := json.Marshal(d1)
	assert.Equal(t, err, nil)
	//使用简单简化数据来解析：
	var dst2 = &Demo4Json{}
	err = json.Unmarshal(data, dst2)
	assert.Equal(t, err, nil)
}
