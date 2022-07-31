// Package githubparse github仓库解析
package githubparse

import (
	"regexp"

	zero "github.com/wdvxdr1123/ZeroBot"

	"github.com/DaydreamCafe/Cocoa/V2/utils/control"
)

// GithubRegex 链接正则表达式
const GithubRegex = `github\.com\/([^(\s|\/)]+\/[^(\s|\/)]+)`

// compiledGithubRegex 编译后的链接正则表达式
var compiledGithubRegex = regexp.MustCompile(GithubRegex)

func init() {
	// 设置插件信息
	metadata := control.Metadata{
		Name:        "github_parse",
		Version:     "1.0.0",
		Description: "获得github仓库简介",
		Author:      "jiangnan777312 / WhitePaper233",
		Usage:       "github链接  得到匹配的仓库简介",
	}

	// 初始化插件
	engine := control.Registe(&metadata, control.EchoAny)

	// 处理github仓库链接
	engine.OnRegex(GithubRegex, zero.OnlyGroup).Handle(
		control.CheckPremissionHandler(handleGithubLink, 5, control.OnlyEchoError),
	)
}
