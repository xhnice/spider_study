package configuration

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

// config
var (
	ConfigFile string
	cfg        Conf
)

// Log -
type Log struct {
	IsStdout bool   `yaml:"isStdout"`
	LogFile  string `yaml:"logFile"`
	Level    uint32 `yaml:"level"`
}

// Conf 配置文件
type Conf struct {
	Log      Log               `yaml:"log"`
	DataPath map[string]string `yaml:"dataPath"`
	DataURI  map[string]string `yaml:"dataUri"`
}

// Load 加载配置文件
func Load() error {
	data, err := ioutil.ReadFile(ConfigFile)
	if err != nil {
		return err
	}
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return err
	}
	return nil
}

// Config 获取配置文件信息
func Config() *Conf {
	return &cfg
}
