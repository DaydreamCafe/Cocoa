/*
@Title        init.go
@Description  bilibili_parse 插件注册
@Author       WhitePaper233 2022.7.2
@Update       WhitePaper233 2022.7.13 
*/
package bilibili_parse

import (
	"regexp"

	zero "github.com/wdvxdr1123/ZeroBot"

	"github.com/DaydreamCafe/Cocoa/V2/utils/plugin"
)

const (
	// 番号正则表达式
	VID_REGEX = `((av|AV)\d+|(bv|BV)1(\d|\w){2}4(\d|\w)1(\d|\w)7(\d|\w){2})`
	// 短链接正则表达式
	SHORT_LINK_REGEX = `(https:\/\/)?b23.tv\/\S{7}`
)

var (
	// 插件元数据
	Metadata plugin.Metadata

	// 编译后的番号正则表达式
	CompiledVIDRegex *regexp.Regexp
	// 编译后的短链接正则表达式
	CompiledShortLinkRegex *regexp.Regexp
)

func init() {
	// 设置插件信息
	Metadata := plugin.Metadata{
		Name:        "bilibili_parse",
		Version:     "1.0.0",
		Description: "bilibili视频解析插件",
		Author: 	"WhitePaper233",
	}
	// 初始化插件
	plugin.Initialization(Metadata)

	// 编译正则表达式
	CompiledVIDRegex = regexp.MustCompile(VID_REGEX)
	CompiledShortLinkRegex = regexp.MustCompile(SHORT_LINK_REGEX)

	// 处理av号或者BV号
	zero.OnRegex(VID_REGEX).Handle(HandleVideoID)

	// 匹配移动端卡片分享信息
	zero.OnMessage().Handle(HandleMobileShare)

	// 匹配短链接
	zero.OnRegex(SHORT_LINK_REGEX).Handle(HandleShortLink)
}
