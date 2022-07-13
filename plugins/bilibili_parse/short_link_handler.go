/*
@Title        short_link_handler.go
@Description  短链接分享处理
@Author       WhitePaper233 2022.7.13
@Update       WhitePaper233 2022.7.13
*/
package bilibili_parse

import (
	"io/ioutil"
	"net/http"
	"strings"

	logger "github.com/sirupsen/logrus"
	zero "github.com/wdvxdr1123/ZeroBot"
)

func HandleShortLink(ctx *zero.Ctx) {
	// 如果是卡片信息, 则跳过
	if strings.HasPrefix(
		ctx.Event.RawMessage,
		`[CQ:json,data={"app":"com.tencent.miniapp_01"&#44;"appID":"100951776"&#44;`,
	) {
		return
	}

	logger.Debugln("匹配短链分享信息成功,MessageId:", ctx.Event.MessageID)
	// 匹配结果
	results := CompiledShortLinkRegex.FindAllStringSubmatch(ctx.MessageString(), -1)

	// 构造请求链接
	var shortLinks []string
	for _, result := range results {
		if strings.HasPrefix(result[0], "https://") {
			// 如果是https链接，则直接使用
			shortLinks = append(shortLinks, result[0])
		} else if strings.HasPrefix(result[0], "http://") {
			// 如果是http链接，则替换为https链接
			shortLinks = append(shortLinks, strings.Replace(result[0], "http://", "https://", 1))
		} else {
			// 如果是相对链接，则拼接成https协议的绝对链接
			shortLinks = append(shortLinks, "https://"+result[0])
		}
	}
	
	// 获取视频长链接
	var fullLinks []string
	for _, shortLink := range shortLinks {
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

		fullLinks = append(fullLinks, fullLink)
	}

	// 匹配结果
	var videoIDs []string
	for _, fullLink := range fullLinks {
		results := CompiledVIDRegex.FindStringSubmatch(fullLink)
		vid := results[0]
		videoIDs = append(videoIDs, vid)
		logger.Debugln("获取视频ID成功:", vid)
	}

	// 获取视频信息
	var videoInfos []VideoInfo
	for _, vid := range videoIDs {
		videoInfo, err := GetVideoInfo(vid)
		if err != nil {
			logger.Errorln("获取视频信息失败:", err)
			continue
		}
		videoInfos = append(videoInfos, videoInfo)
	}

	// 发送视频信息
	for _, videoInfo := range videoInfos {
		videoInfo.Send(ctx)
	}
}
