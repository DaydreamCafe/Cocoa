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

var botHelp = `
`

func HandleHelp(ctx *zero.Ctx) {
	fset := flag.FlagSet{}

	// plugin
	var	plugin string
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
		ctx.Send(GetPluginUsage(plugin))
	case bot:
		// 显示bot的帮助
		ctx.Send(botHelp)
	case help:
		// 显示help指令的帮助
		ctx.Send(HelpCommandHelp)
	default:
		// 显示help指令的帮助
		ctx.Send(HelpCommandHelp)
	}
}

func GetPluginUsage(name string) string {
	db, err := conn.GetDB()
	if err != nil {
		logger.Error("获取数据库连接失败:", err)
		return "未找到该插件的帮助信息"
	}

	plugin_metadata := model.PluginModel{}
	db.Where("name = ?", name).First(&plugin_metadata)
	if plugin_metadata.Usage == "" {
		return "未找到该插件的帮助信息"
	}
	return fmt.Sprint(plugin_metadata.Usage)
}