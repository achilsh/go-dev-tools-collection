package config

import (
	"fmt"
	"io/ioutil"

	"github.com/go-yaml/yaml"
)

var (
	Global = new(GlobalConfig)
)

type DBCfg struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	DBName   string `yaml:"dbname"`
	UserName string `yaml:"username"`
	Pwd      string `yaml:"pwd"`
}
type LoggerCfg struct {
	LogPath string `yaml:"logPath"`
}
type RedisCfg struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
	DataBase int    `yaml:"database"`
}
type HttpServerCfg struct {
	Port int `yaml:"port"`
}

type ApiSecretKey struct {
	APISignKey      string `yaml:"apiSignKey"`      //API_SIGN_KEY, sign used
	AESKeyToken     string `yaml:"aesKeyToken"`     // AES_KEY used token
	AESIVToken      string `yaml:"aesIVToken"`      //AES_IV used token
	RefreshTokenKey string `yaml:"refreshTokenKey"` // REFRESH_TOKEN_KEY
}

// GlobalConfig 全局配置
type GlobalConfig struct {
	DB          *DBCfg               `yaml:"db"`
	LoggerItem  *LoggerCfg           `yaml:"logger"`
	RedisItem   *RedisCfg            `yaml:"redis"`
	Server      *HttpServerCfg       `yaml:"server"`
	SecretKey   *ApiSecretKey        `yaml:"secretKey"`
	WelWordLang *WelcomeWordLanguage `yaml:"welWordLang"`
}

// GetGlobalConfig 获取全局配置
func GetGlobalConfig() *GlobalConfig {
	return Global
}

func ParseYamlFile(filename string) (*GlobalConfig, error) {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("ReadFile: %v\n", err)
		return nil, err
	}
	// 初始化 ServiceCfg
	Global = &GlobalConfig{}
	err = yaml.Unmarshal(buf, Global)
	if err != nil {
		fmt.Printf("Unmarshal: %v\n", err)

		return Global, err
	}
	fmt.Printf("ServiceCfg: %+v\n", Global)
	return Global, nil
}
