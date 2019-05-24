package utils

import (
	"fmt"
	"sync_study/configuration"
	"sync_study/logger"
	"testing"
)

func init() {

	configuration.ConfigFile = `./../config.yaml`

	configuration.Load()

	cfg := configuration.Config()
	err := logger.Init(cfg.Log.IsStdout, "../"+cfg.Log.LogFile, cfg.Log.Level)
	if err != nil {
		panic(fmt.Errorf("日志信息配置失败~%s", err))
	}
}

func TestFile(t *testing.T) {

	var (
		filePath string
		err      error
	)

	filePath = "../logs"
	if !FileExists(filePath) {
		t.Errorf("%s 目录存在", filePath)
	}

	filePath = "../logs/empty.log"
	if FileExists(filePath) {
		t.Errorf("%s 目录不存在", filePath)
	}

	filePath = "../logs"

	_, err = OpenFile(filePath)
	if err == nil {
		t.Errorf("打开了一个目录[%s]", filePath)
	} else {
		t.Logf("打开目录 [%s] 出错[%s]", filePath, err)
	}

	filePath = "../logs/test.file.log"
	_, err = OpenFile(filePath)
	if err != nil {
		t.Errorf("打开文件[%s]出错了[%s]", filePath, err)
	}
}
