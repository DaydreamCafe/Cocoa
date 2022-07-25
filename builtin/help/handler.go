package help

import (
	"flag"
	"fmt"

	logger "github.com/sirupsen/logrus"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/extension/shell"

	"github.com/DaydreamCafe/Cocoa/V2/src/conn"
	"github.com/DaydreamCafe/Cocoa/V2/src/model"
)

const botHelp = `Cocoa 是使用 Go语言 和 ZeroBot框架 编写的一款小巧的QQ机器人; Cocoa可以与任何支持 OneBot标准 的机器人框架/平台进行交互.

Cocoa 的初衷是构建一个高可用且小巧的QQ机器人后端; 并提供一套实用且足够精简的插件, 在尽可能少向QQ群里制造信息的同时最大化服务于您的QQ群。但是, 您仍然可以通过安装其他支持 Cocoa 的插件丰富其功能.

NOTICE: Cocoa 仍处于开发早期阶段, 很多功能并没有稳定, 不建议您立刻将其用于生产环境.

https://github.com/DaydreamCafe/Cocoa`

// handleHelp 处理帮助信息
func handleHelp(ctx *zero.Ctx) {
	fset := flag.FlagSet{}

	// plugin
	var plugin string
	fset.StringVar(&plugin, "p", "", "显示插件的帮助")
	fset.StringVar(&plugin, "plugin", "", "显示插件的帮助")

	// bot
	var bot bool
	fset.BoolVar(&bot, "b", false, "显示bot的帮助")
	fset.BoolVar(&bot, "bot", false, "显示bot的帮助")

	// help
	var help bool
	fset.BoolVar(&help, "h", false, "显示help指令的帮助")
	fset.BoolVar(&help, "help", false, "显示help指令的帮助")

	args := shell.Parse(ctx.State["args"].(string))
	err := fset.Parse(args)
	if err != nil {
		return
	}

	switch {
	case plugin != "":
		// 显示插件的帮助
		ctx.Send(getPluginUsage(plugin))
	case bot:
		// 显示bot的帮助
		ctx.Send(botHelp)
	case help:
		// 显示help指令的帮助
		ctx.Send(usage)
	default:
		// 显示help指令的帮助
		ctx.Send(usage)
	}
}

// getPluginUsage 获取插件的帮助
func getPluginUsage(name string) string {
	db, err := conn.GetDB()
	if err != nil {
		logger.Error("获取数据库连接失败:", err)
		return "未找到该插件的帮助信息"
	}

	sqlDB, err := db.DB()
	if err != nil {
		logger.Error("获取数据库连接失败:", err)
		return "未找到该插件的帮助信息"
	}
	defer sqlDB.Close()

	plugin_metadata := model.Plugin{}
	db.Where("name = ?", name).First(&plugin_metadata)
	if plugin_metadata.Usage == "" {
		return "未找到该插件的帮助信息"
	}
	return fmt.Sprint(plugin_metadata.Usage)
}
