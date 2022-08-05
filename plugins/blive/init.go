// package blive B站直播搜索
package blive

import (
	"fmt"
	"regexp"

	"github.com/DaydreamCafe/Cocoa/V2/utils/control"
	zero "github.com/wdvxdr1123/ZeroBot"
)

const (
	// CommandRegex 直播间搜索正则表达式
	CommandRegex = `^blive \S+`
)

var (
	// compiledCommandRegex 编译后的直播间搜索正则表达式
	compiledCommandRegex *regexp.Regexp = regexp.MustCompile(CommandRegex)
)

func init() {
	// 设置插件信息
	metadata := control.Metadata{
		Name:        "blive",
		Version:     "1.0.0",
		Description: "Bilibili直播搜索",
		Author:      "jiangnan777312 / WhitePaper233",
		Usage:       `-提供B站直播搜索`,
	}
	// 初始化插件
	engine := control.Registe(&metadata, control.EchoAny)

	// 处理直播间搜索
	engine.OnRegex(CommandRegex, zero.OnlyGroup).Handle(
		control.CheckPremissionHandler(handleBlive, 5, control.OnlyEchoError),
	)
}

// 格式化数字
func formatDigit(digit int) string {
	// 当数字大于9999时, 显示为"x.x万"
	// 当数字大于99999999时, 显示为"x.x亿"
	if digit > 99999999 {
		return fmt.Sprintf("%.1f亿", float64(digit)/100000000)
	} else if digit > 99999 {
		return fmt.Sprintf("%.1f万", float64(digit)/10000)
	}
	return fmt.Sprintf("%d", digit)
}
