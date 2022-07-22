// Package charreverser 英文字符翻转
package charreverser

import (
	zero "github.com/wdvxdr1123/ZeroBot"

	"github.com/DaydreamCafe/Cocoa/V2/utils/control"
)

func init() {
	// 设置插件信息
	metadata := control.Metadata{
		Name:        "char_reverser",
		Version:     "1.0.0",
		Description: "翻转英文字符",
		Author:      "WhitePaper233",
		Usage:       "CharReverser插件\n- 翻转 <英文字符串>  得到一个翻转的英文字符",
	}
	// 初始化插件
	control.Registe(&metadata)

	// 处理翻转命令
	zero.OnRegex(`翻转( )+[A-z]{1}([A-z]|\s)+[A-z]{1}`).SetBlock(true).Handle(handleReverse)
}
