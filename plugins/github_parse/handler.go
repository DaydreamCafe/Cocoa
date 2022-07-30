// Package githubparse github仓库解析
package githubparse

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	logger "github.com/sirupsen/logrus"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

// apiURL API地址
const apiURL = "https://api.github.com/repos"

// LicenseData 许可证信息结构体
type LicenseData struct {
	Name string `json:"name"`
}

// APIResp GithubAPI返回结果结构体
type APIResp struct {
	Name        string      `json:"name"`
	FullName    string      `json:"full_name"`
	URL         string      `json:"html_url"`
	Description string      `json:"description"`
	License     LicenseData `json:"license"`
	Forks       uint64      `json:"forks"`
	Stars       uint64      `json:"stargazers_count"`
	Watchers    uint64      `json:"watchers"`
}

// handleGithubLink GithubLinkHandler处理Github链接
func handleGithubLink(ctx *zero.Ctx) {
	// 提取链接匹配结果
	links := compiledGithubRegex.FindAllStringSubmatch(ctx.ExtractPlainText(), -1)
	var repos = make([]string, len(links))
	for index, link := range links {
		repos[index] = link[1]
	}

	// 构造请求URL
	var reqURLs = make([]string, len(links))
	for index, repo := range repos {
		var reqURLBuilder strings.Builder

		reqURLBuilder.WriteString(apiURL)
		reqURLBuilder.WriteRune('/')
		reqURLBuilder.WriteString(repo)

		reqURLs[index] = reqURLBuilder.String()
	}

	//从API获取仓库信息, 并解析为结构体
	var respInfos = make([]APIResp, len(reqURLs))
	for index, reqURL := range reqURLs {
		// 调用API获取信息
		response, err := http.Get(reqURL)
		if err != nil {
			logger.Errorln("请求失败:", err)
			continue
		}
		// 当状态码不为200时
		if response.StatusCode != 200 {
			logger.Debugln(response.StatusCode)
			logger.Warningln("仓库不存在:", repos[index])
			continue
		}
		defer response.Body.Close()

		// 将请求结果JSON解析为APIResp结构体
		var respInfo APIResp
		err = json.NewDecoder(response.Body).Decode(&respInfo)
		if err != nil {
			logger.Errorln("JSON解析失败:", err)
			continue
		}

		respInfos[index] = respInfo
	}

	// 格式化回复字符串
	var replyStr = make([]string, len(respInfos))
	for index, respInfo := range respInfos {
		// 上一步可能出现请求失败或请求仓库不存在的情况, 在下面过滤掉, 直接设为空字符串
		if respInfo.Name == "" {
			replyStr[index] = ""
			continue
		}
		var replyBuilder strings.Builder
		replyBuilder.WriteString(respInfo.Name)
		replyBuilder.WriteRune('\n')
		replyBuilder.WriteString(respInfo.FullName)
		replyBuilder.WriteRune('\n')
		replyBuilder.WriteString("--------------------\n")
		replyBuilder.WriteString("stars: ")
		replyBuilder.WriteString(strconv.FormatUint(respInfo.Stars, 10))
		replyBuilder.WriteString("   ")
		replyBuilder.WriteString("forks: ")
		replyBuilder.WriteString(strconv.FormatUint(respInfo.Forks, 10))
		replyBuilder.WriteString("   ")
		replyBuilder.WriteString("watchers: ")
		replyBuilder.WriteString(strconv.FormatUint(respInfo.Watchers, 10))
		replyBuilder.WriteString("\n\n")
		replyBuilder.WriteString(respInfo.Description)
		replyBuilder.WriteString("\n\n")
		replyBuilder.WriteString("license: ")
		replyBuilder.WriteString(respInfo.License.Name)
		replyBuilder.WriteString("\n\n")
		replyBuilder.WriteString(respInfo.URL)

		replyStr[index] = replyBuilder.String()
	}

	// 发送信息
	for _, reply := range replyStr {
		// 过滤掉空字符串
		if reply != "" {
			ctx.SendChain(message.Text(reply))
		}
	}
}
