package help

import (
	zero "github.com/wdvxdr1123/ZeroBot"

	"github.com/DaydreamCafe/Cocoa/V2/utils/control"
)


var HelpCommandHelp = `帮助插件
help -<option> <arg>
	-p  -plugin <插件名>  显示插件的帮助
	-b  -bot  显示bot的帮助
	-h  -help  显示help指令的帮助
`

func init() {
	// 设置插件信息
	Metadata := control.Metadata{
		Name:        "help",
		Version:     "1.0.0",
		Description: "帮助插件",
		Author:      "WhitePaper233",
		Usage:       HelpCommandHelp,
	}
	// 初始化插件
	control.Registe(&Metadata)

	// 处理help命令
	var helpCommands = []string{
		"help",
		"帮助",
	}
	zero.OnCommandGroup(helpCommands).Handle(HandleHelp)
}
