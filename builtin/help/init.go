package help

import (
	"github.com/DaydreamCafe/Cocoa/V2/utils/control"
)

const usage = `帮助插件
- help -<option> <arg>
- 帮助 -<option> <arg>
	-p  -plugin <插件名>  显示插件的帮助
	-b  -bot  显示bot的帮助
	-h  -help  显示help指令的帮助
`

func init() {
	// 设置插件信息
	metadata := control.Metadata{
		Name:        "help",
		Version:     "1.0.0",
		Description: "帮助插件",
		Author:      "WhitePaper233",
		Usage:       usage,
		Buitlin:     true,
	}
	// 初始化插件
	engine := control.Registe(&metadata)

	// 处理help命令
	var helpCommands = []string{
		"help",
		"帮助",
	}
	engine.OnCommandGroup(helpCommands).Handle(handleHelp).SetBlock(true)
}
