// Package init Bot初始化相关代码
package init

import (
	logger "github.com/sirupsen/logrus"
	loggerFormatter "github.com/t-tomalak/logrus-easy-formatter"
	"gopkg.in/yaml.v3"

	"github.com/DaydreamCafe/Cocoa/V2/src/io"
)

/*
CONFIG CONFIG结构体, 用于储存配置文件中的DEBUG选项
因为logger初始化需在全部读取配置文件之前, 所以写了一个简化的配置文件读取操作
*/
type CONFIG struct {
	DEBUG bool `yaml:"DEBUG"`
}

func init() {
	// 设置logger基础样式
	logger.SetFormatter(&loggerFormatter.Formatter{
		TimestampFormat: "2006-01-02 15:04:05",
		LogFormat:       "[Cocoa][%time%][%lvl%]: %msg%\n",
	})

	// 从文件中读取配置文件
	configFile, err := io.ReadConfig()
	if err != nil {
		logger.Panic(err)
		panic(err)
	}

	var config CONFIG
	// 解析配置文件
	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		logger.Error(err)
		config.DEBUG = false
	}
	// 设置DEBUG是否为debug模式的logger
	if config.DEBUG {
		logger.SetLevel(logger.DebugLevel)
		logger.Debugln("正在以 DEBUG 模式运行...")
	} else {
		logger.SetLevel(logger.InfoLevel)
	}
}
