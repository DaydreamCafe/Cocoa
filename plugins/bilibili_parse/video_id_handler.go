// Package bilibiliparse B站分享解析
package bilibiliparse

import (
	logger "github.com/sirupsen/logrus"
	zero "github.com/wdvxdr1123/ZeroBot"
)

// handleVideoID 视频ID Handler
func handleVideoID(ctx *zero.Ctx) {
	// 匹配结果
	results := compiledVIDRegex.FindAllStringSubmatch(ctx.MessageString(), -1)

	// 循环匹配结果
	var VIDs = make([]string, len(results))

	for index, result := range results {
		// 获取av号或者BV号
		VIDs[index] = result[0]
		logger.Debugln("匹配视频ID成功:", result[0], ", MessageId:", ctx.Event.MessageID)
	}

	for _, vid := range VIDs {
		// 获取视频信息
		videoInfo, err := getVideoInfo(vid)
		if err != nil {
			logger.Errorln("获取视频信息失败:", err)
			continue
		}
		// 发送视频信息
		videoInfo.send(ctx)
	}
}
