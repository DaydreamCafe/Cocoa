// Package coser coser图片
package coser

import (
	zero "github.com/wdvxdr1123/ZeroBot"

	"github.com/DaydreamCafe/Cocoa/V2/utils/control"
)

func init() {
	// 设置插件信息
	Metadata := control.Metadata{
		Name:        "coser",
		Version:     "1.0.0",
		Description: "三次元小姐姐",
		Author:      "WhitePaper233",
		Usage: `Coser插件
		coser  得到一张coser的图片`,
	}
	// 初始化插件
	control.Registe(&Metadata)

	// 处理coser命令
	zero.OnFullMatch("coser", zero.OnlyGroup).SetBlock(true).Handle(HandleCoser)
}
