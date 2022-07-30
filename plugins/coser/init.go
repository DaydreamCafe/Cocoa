// Package coser coser图片
package coser

import (
	zero "github.com/wdvxdr1123/ZeroBot"

	"github.com/DaydreamCafe/Cocoa/V2/utils/control"
)

func init() {
	// 设置插件信息
	metadata := control.Metadata{
		Name:        "coser",
		Version:     "1.0.0",
		Description: "三次元小姐姐",
		Author:      "WhitePaper233",
		Usage:       "Coser插件\n-coser  得到一张coser的图片",
	}
	// 初始化插件
	engine := control.Registe(&metadata, control.EchoAny)

	// 处理coser命令
	engine.OnFullMatch("coser", zero.OnlyGroup).SetBlock(true).Handle(
		control.CheckPremissionHandler(handleCoser, 5, control.EchoAny),
	)
}
