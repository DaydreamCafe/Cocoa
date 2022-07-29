// Package pluginmanager 插件管理器
package pluginmanager

import (
	zero "github.com/wdvxdr1123/ZeroBot"

	"github.com/DaydreamCafe/Cocoa/V2/utils/control"
)

const usage = `插件管理器
- plugin -<option> <arg>
- 插件 -<option> <arg>
	-l  -list    显示插件列表
	-b  -ban     禁用插件
	-u  -unban   解禁插件
	-e  -enable  在当前群启用插件
	-d  -disable 在当前群禁用插件
	-h  -help    显示此帮助`

func init() {
	// 设置插件信息
	metadata := control.Metadata{
		Name:        "plugin_manager",
		Version:     "1.0.0",
		Description: "插件管理器",
		Author:      "WhitePaper233",
		Usage:       usage,
		Buitlin:     true,
	}
	// 初始化插件
	engine := control.Registe(&metadata, control.EchoAny)

	// 更新插件管理表
	initDatabase()

	// 处理插件管理命令
	var pluginCommands = []string{
		"plugin",
		"插件",
	}
	engine.
	OnCommandGroup(pluginCommands, zero.OnlyGroup).Handle(handlePlugin).SetBlock(true)
}
