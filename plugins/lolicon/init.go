// Package lolicon lolicon图片
package lolicon

import (
	zero "github.com/wdvxdr1123/ZeroBot"

	"github.com/DaydreamCafe/Cocoa/V2/utils/control"
)

func init() {
	// 设置插件信息
	metadata := control.Metadata{
		Name:        "lolicon",
		Version:     "1.0.0",
		Description: "二次元涩图",
		Author:      "jiangnan777312 / WhitePaper233",
		Usage:       "Lolicon插件\n-涩图|色图  得到一张loli的图片",
	}
	// 初始化插件
	engine := control.Registe(&metadata, control.EchoAny)

	// 处理涩图命令
	commandGroup := []string{
		"涩图",
		"色图",
	}
	engine.OnPrefixGroup(commandGroup, zero.OnlyGroup).SetBlock(true).Handle(
		control.CheckPremissionHandler(handleLoli, 5, control.EchoAny),
	)
}
