package example

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/achilsh/go-dev-tools-collection/base-lib/config/file"
	logger "github.com/achilsh/go-dev-tools-collection/base-lib/log"
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
	
}
