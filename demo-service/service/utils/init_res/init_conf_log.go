package init_res

import (
	"fmt"

	"github.com/achilsh/go-dev-tools-collection/demo-service/service/utils/config"

	logger "github.com/achilsh/go-dev-tools-collection/base-lib/log"
)

func InitConfigLog(aiConfigFile *string) error {
	if _, err := config.ParseYamlFile(*aiConfigFile); err != nil {
		fmt.Println("init config fail, err: ", err)
		return fmt.Errorf("parse config fail, err: %v", err)
	}
	if aiConfigFile == nil {
		fmt.Println("get demo config path is nil")
		return fmt.Errorf("config is nil")
	}
	fmt.Println("demo server config path: ", *aiConfigFile)

	fmt.Println("demo config file: ", config.GetGlobalConfig())
	if err := logger.Init(*aiConfigFile, config.GetGlobalConfig().LoggerItem.LogPath); err != nil {
		panic(err)
	}

	logger.Debugf("config: %+v", config.GetGlobalConfig().DB)
	return nil
}
