// Package github github仓库解析
package githubparse

import (
	"encoding/json"
	"net/http"
	"strings"

	logger "github.com/sirupsen/logrus"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

const (
	// apiURL API地址
	apiURL = "https://api.github.com/repos"
)

// LicenseData结构
type LicenseData struct {
	Key string `json:"key`
}

// ParentData结构
type ParentData struct {
	Name     string `json:"name"`
	FullName string `json:"full_name"`
	URL      string `json:"html_url`
}

// SourceData结构
type SourceData struct {
	Name     string `json:"name"`
	FullName string `json:"full_name"`
	URL      string `json:"html_url`
}

// APIResp GithubAPI返回结果结构体
type APIResp struct {
	Name        string        `json:"name"`
	FullName    string        `json:"full_name`
	URL         string        `json:"html_url`
	Description string        `json:"description`
	License     []LicenseData `json:"license`
	Forks       uint64        `json:"forks"`
	Watchers    uint64        `json:"watchers"`
	Parent      []ParentData  `json:"parent"`
	Source      []SourceData  `json:"source"`
}

// handleGithubLink GithubLinkHandler
func handleGithubLink(ctx *zero.Ctx) {
	// 请求URL
	var reqURL string = apiURL

	// 提取链接匹配结果
	msg := ctx.ExtractPlainText()
	params := compiledGithubRegex.FindStringSubmatch(msg)
	res := params[1]

	// 构造请求URL
	var reqURLBuilder strings.Builder
	reqURLBuilder.WriteString(apiURL)
	reqURLBuilder.WriteRune('/')
	reqURLBuilder.WriteString(res)
	reqURL = reqURLBuilder.String()
	logger.Debugln(reqURL)

	//从API获取仓库信息
	response, err := http.Get(reqURL)
	if err != nil {
		logger.Errorln("请求失败:", err)
		return
	}
	defer response.Body.Close()

	// 将请求结果JSON解析为APIResp结构体
	var resp APIResp
	err = json.NewDecoder(response.Body).Decode(&resp)
	if err != nil {
		logger.Errorln("JSON解析失败:", err)
		return
	}

	// 发送信息
	ctx.SendChain(
		message.Reply(ctx.Event.MessageID),
		message.Text(reply),
	)
}
