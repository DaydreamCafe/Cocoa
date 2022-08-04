// Package bilibiliparse B站分享解析
package bilibiliparse

import (
	"regexp"

	zero "github.com/wdvxdr1123/ZeroBot"

	"github.com/DaydreamCafe/Cocoa/V2/utils/control"
)

const (
	// VIDRegex 番号正则表达式
	VIDRegex = `((av|AV)\d+|(bv|BV)1(\d|\w){2}4(\d|\w)1(\d|\w)7(\d|\w){2})`
	// ShortLinkRegex 短链接正则表达式
	ShortLinkRegex = `(https:\/\/)?b23.tv\/\S{7}`
)

var (
	// compiledVIDRegex 编译后的番号正则表达式
	compiledVIDRegex *regexp.Regexp = regexp.MustCompile(VIDRegex)
	// compiledShortLinkRegex 编译后的短链接正则表达式
	compiledShortLinkRegex *regexp.Regexp = regexp.MustCompile(ShortLinkRegex)
)

func init() {
	// 设置插件信息
	metadata := control.Metadata{
		Name:        "bilibili_parse",
		Version:     "1.0.0",
		Description: "Bilibili视频解析插件",
		Author:      "WhitePaper233",
		Usage:       `-发送任意形式的B站分享链接、番号及移动端分享卡片, 将自动解析出视频信息`,
	}
	// 初始化插件
	engine := control.Registe(&metadata, control.EchoAny)

	// 处理av号或者BV号
	engine.OnRegex(VIDRegex, zero.OnlyGroup).Handle(
		control.CheckPremissionHandler(handleVideoID, 5, control.OnlyEchoError),
	)

	// 匹配短链接和卡片信息
	engine.OnRegex(ShortLinkRegex, zero.OnlyGroup).Handle(
		control.CheckPremissionHandler(handleShortLink, 5, control.OnlyEchoError),
	)
}
