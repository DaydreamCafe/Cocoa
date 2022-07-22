// Package lolicon lolicon图片
package lolicon

import (
	zero "github.com/wdvxdr1123/ZeroBot"

	"github.com/DaydreamCafe/Cocoa/V2/utils/control"
)

func init() {
	// 设置插件信息
	Metadata := control.Metadata{
		Name:        "lolicon",
		Version:     "1.0.0",
		Description: "二次元萝莉",
		Author:      "jiangnan777312",
		Usage: `Lolicon插件
		get_loli  得到一张loli的图片`,
	}
	// 初始化插件
	control.Registe(&Metadata)

	// 处理get_loli命令
	commandGroup := []string{
		"涩图",
		"色图",
	}
	zero.OnFullMatchGroup(commandGroup, zero.OnlyGroup).SetBlock(true).Handle(HandleLoli)
}
