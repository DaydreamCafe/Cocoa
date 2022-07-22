// Package bilibiliparse B站分享解析
package bilibiliparse

import (
	"io/ioutil"
	"net/http"
	"strings"

	logger "github.com/sirupsen/logrus"
	zero "github.com/wdvxdr1123/ZeroBot"
)

// handleShortLink 短链接handler
func handleShortLink(ctx *zero.Ctx) {
	// 如果是卡片信息, 则跳过
	if strings.HasPrefix(
		ctx.Event.RawMessage,
		`[CQ:json,data={"app":"com.tencent.miniapp_01"&#44;`,
	) {
		return
	}

	logger.Debugln("匹配短链分享信息成功,MessageId:", ctx.Event.MessageID)
	// 匹配结果
	results := compiledShortLinkRegex.FindAllStringSubmatch(ctx.MessageString(), -1)

	// 构造请求链接
	var shortLinks = make([]string, len(results))

	for index, result := range results {
		if strings.HasPrefix(result[0], "https://") {
			// 如果是https链接，则直接使用
			shortLinks[index] = result[0]
		} else if strings.HasPrefix(result[0], "http://") {
			// 如果是http链接，则替换为https链接
			shortLinks[index] = strings.Replace(result[0], "http://", "https://", 1)
		} else {
			// 如果是相对链接，则拼接成https协议的绝对链接
			shortLinks[index] = "https://" + result[0]
		}
	}

	// 获取视频长链接
	var fullLinks = make([]string, len(shortLinks))

	for index, shortLink := range shortLinks {
		request, err := http.NewRequest("GET", shortLink, nil)
		if err != nil {
			logger.Errorln("创建请求失败:", err)
			continue
		}

		response, err := http.DefaultClient.Do(request)
		if err != nil {
			logger.Errorln("请求失败:", err)
			continue
		}
		defer response.Body.Close()

		// 解析响应为字符串
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			logger.Errorln("读取响应失败:", err)
			continue
		}
		// 抽取重定向链接
		redictPage := string(body)

		var fullLink string
		fullLink = strings.TrimPrefix(redictPage, "<a href=\"")
		fullLink = strings.TrimSuffix(fullLink, "\">Found</a>.")

		fullLinks[index] = fullLink
	}

	// 匹配结果
	var videoIDs = make([]string, len(fullLinks))

	for index, fullLink := range fullLinks {
		results := compiledVIDRegex.FindStringSubmatch(fullLink)
		if len(results) == 0 {
			continue
		}

		vid := results[0]
		videoIDs[index] = vid
		logger.Debugln("获取视频ID成功:", vid)
	}

	// 获取视频信息
	var videoInfos = make([]VideoInfo, len(videoIDs))

	for index, vid := range videoIDs {
		if vid == "" {
			continue
		}

		videoInfo, err := getVideoInfo(vid)
		if err != nil {
			logger.Errorln("获取视频信息失败:", err)
			continue
		}
		videoInfos[index] = videoInfo
	}

	// 发送视频信息
	for _, videoInfo := range videoInfos {
		if videoInfo.BVID == "" {
			continue
		}

		videoInfo.send(ctx)
	}
}
