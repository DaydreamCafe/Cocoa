// Package usermanager 用户管理相关代码
package usermanager

import (
	zero "github.com/wdvxdr1123/ZeroBot"

	"github.com/DaydreamCafe/Cocoa/V2/utils/control"
)

const usage = `用户管理器
- user -<option> <arg>
- 用户 -<option> <arg>
	-u   -user   <ID>     选择操作用户
	-s   -set    <level>  设置用户等级
	-b   -ban    <time>   禁用用户
	-p   -pardon          解禁用户
	-r   -reset           重置用户
	-h   -help            显示此帮助`

func init() {
	// 设置插件信息
	metadata := control.Metadata{
		Name:        "user_manager",
		Version:     "1.0.0",
		Description: "用户管理器",
		Author:      "WhitePaper233",
		Usage:       usage,
		Buitlin:     true,
	}
	// 初始化插件
	engine := control.Registe(&metadata, control.EchoAny)

	// 更新用户相关数据表
	initDatabase()

	// 处理插件管理命令
	var userCommands = []string{
		"user",
		"用户",
	}
	engine.OnCommandGroup(userCommands, zero.OnlyGroup).Handle(handleUser).SetBlock(true)
}
