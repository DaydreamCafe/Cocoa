package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"

	log "github.com/sirupsen/logrus"
	easy "github.com/t-tomalak/logrus-easy-formatter"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/driver"
	yaml "gopkg.in/yaml.v2"

	// 加载插件
	_ "github.com/DaydreamCafe/Cocoa/plugins/nickname"
)

// 配置文件结构体
type Config struct {
	SuperUsers    []string `yaml:"SuperUsers"`
	NickNames     []string `yaml:"NickNames"`
	CommandPrefix string   `yaml:"CommandPrefix"`
	DEBUG         bool     `yaml:"DEBUG"`
	Server        struct {
		Address string `yaml:"Address"`
		Port    int64    `yaml:"Port"`
		Token   string `yaml:"Token"`
	} `yaml:"Server"`
}

var config Config

func init() {
	// 初始化logger
	log.SetFormatter(&easy.Formatter{
		TimestampFormat: "2006-01-02 15:04:05",
		LogFormat:       "[zero][%time%][%lvl%]: %msg% \n",
	})

	// 初始化配置文件
	configFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Error(err)
	}
	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		log.Error(err)
	}

	// DEBUG模式
	if config.DEBUG {
		log.SetLevel(log.DebugLevel)
		log.Debug("正在以DEBUG模式运行程序")
	}
}

func main() {
	// 初始化框架
	zeroCondfig := zero.Config{
		NickName:      config.NickNames,
		CommandPrefix: config.CommandPrefix,
		SuperUsers:    config.SuperUsers,
		Driver: []zero.Driver{driver.NewWebSocketClient(
			fmt.Sprintf("ws://%s:%d", config.Server.Address, config.Server.Port),
			config.Server.Token,
		),
		}}
	// 运行插件
	zero.Run(zeroCondfig)

	// 捕获Ctrl C退出程序
	signalChan := make(chan os.Signal, 1)
	cleanupDone := make(chan bool)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		for range signalChan {
			fmt.Printf("\n停止服务...\n")
			cleanUp()
			cleanupDone <- true
		}
	}()
	<-cleanupDone
}

func cleanUp() {
	// 先留着罢，以后需要的时候再写awa
}
