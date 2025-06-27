package gobuildopt

import (
	logger "github.com/achilsh/go-dev-tools-collection/base-lib/log"
)

var (
	//Version 项目版本信息
	Version = ""
	//GoVersion Go版本信息
	GoVersion = ""
	//GitCommit git提交commmit id
	GitCommit = ""
	//BuildTime 构建时间
	BuildTime = ""
)

func ShowVersion() {
	logger.Infof("Version: %s", Version)
	logger.Infof("Go Version: %s", GoVersion)
	logger.Infof("Git Commit: %s", GitCommit)
	logger.Infof("Build Time: %s", BuildTime)
}
