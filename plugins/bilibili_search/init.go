// Package bilibilisearch B站综合搜索
package bilibilisearch

import (
	"fmt"
	"regexp"

	"github.com/DaydreamCafe/Cocoa/V2/utils/control"
	zero "github.com/wdvxdr1123/ZeroBot"
)

const (
	// LiveSearchRegex 直播间搜索正则表达式
	LiveSearchRegex = `^bililive ((\S+)(&(\S+))?)+`
)

var (
	// compiledLiveSearchRegex 编译后的直播间搜索正则表达式
	compiledLiveSearchRegex *regexp.Regexp = regexp.MustCompile(LiveSearchRegex)
)

func init() {
	// 设置插件信息
	metadata := control.Metadata{
		Name:        "bilibili_search",
		Version:     "1.0.0",
		Description: "Bilibili综合搜索插件",
		Author:      "jiangnan777312",
		Usage: `-支持以下B站搜索功能: 
				-直播间搜索`,
	}
	// 初始化插件
	engine := control.Registe(&metadata, control.EchoAny)

	// 处理直播间搜索
	engine.OnRegex(LiveSearchRegex, zero.OnlyGroup).Handle(
		control.CheckPremissionHandler(handleLiveSearch, 5, control.OnlyEchoError),
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
