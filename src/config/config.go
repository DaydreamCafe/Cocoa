// Package config 配置文件操作
package config

import (
	logger "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"

	"github.com/DaydreamCafe/Cocoa/V2/src/io"
)

// Config Config类, 用于读取和操作配置文件
type Config struct {
	SuperUsers    []int64  `yaml:"SuperUsers"`
	CommandPrefix string   `yaml:"CommandPrefix"`
	NickNames     []string `yaml:"NickNames"`
	DefaultLevel  int      `yaml:"DefaultLevel"`

	Server struct {
		Address string `yaml:"Address"`
		Port    int    `yaml:"Port"`
		Token   string `yaml:"Token"`
	} `yaml:"Server"`

	Database struct {
		Address      string `yaml:"Address"`
		Port         int    `yaml:"Port"`
		User         string `yaml:"User"`
		Password     string `yaml:"Password"`
		DatabaseName string `yaml:"DatabaseName"`
	} `yaml:"Database"`

	DEBUG bool `yaml:"DEBUG"`
}

// Load 加载配置文件
func (cfg *Config) Load() {
	configFile, err := io.ReadConfig()
	if err != nil {
		logger.Panic(err)
		panic(err)
	}

	err = yaml.Unmarshal(configFile, &cfg)
	if err != nil {
		logger.Panic(err)
		panic(err)
	}
}
