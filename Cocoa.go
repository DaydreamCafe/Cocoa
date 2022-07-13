package main

import (
	// internal packages
	"fmt"
	"os"
	"os/signal"

	// external packages
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/driver"

	// custom packages
	"github.com/DaydreamCafe/Cocoa/V2/src/config"
	_ "github.com/DaydreamCafe/Cocoa/V2/src/logger"

	// plugin packages
	_ "github.com/DaydreamCafe/Cocoa/V2/plugins/bilibili_parse"
)

var Config config.CONFIG
var zeroConfig zero.Config

// 初始化函数
func init() {
	// 加载配置文件
	Config = Config.Load()

	// 初始化ZeroBot配置
	zeroConfig = zero.Config{
		NickName:      Config.NickNames,
		CommandPrefix: Config.CommandPrefix,
		SuperUsers:    Config.SuperUsers,
		Driver: []zero.Driver{driver.NewWebSocketClient(
			fmt.Sprintf("ws://%s:%d", Config.Server.Address, Config.Server.Port),
			Config.Server.Token,
		),
		},
	}
}

// 入口函数
func main() {
	// 初始化并运行ZeroBot
	zero.Run(zeroConfig)

	// 处理Ctrl+C退出信号
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

// 清理函数
func cleanUp() {

}
