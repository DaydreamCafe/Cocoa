// Package bilibiliparse B站分享解析
package bilibiliparse

import (
	logger "github.com/sirupsen/logrus"
	zero "github.com/wdvxdr1123/ZeroBot"
)

// HandleVideoID 视频ID Handler
func HandleVideoID(ctx *zero.Ctx) {
	// 匹配结果
	results := CompiledVIDRegex.FindAllStringSubmatch(ctx.MessageString(), -1)

	// 循环匹配结果
	var VID = make([]string, len(results))

	for index, result := range results {
		// 获取av号或者BV号
		VID[index] = result[0]
		logger.Debugln("匹配视频ID成功:", result[0], ", MessageId:", ctx.Event.MessageID)
	}

	for _, vid := range VID {
		// 获取视频信息
		videoInfo, err := GetVideoInfo(vid)
		if err != nil {
			logger.Errorln("获取视频信息失败:", err)
			continue
		}
		// 发送视频信息
		videoInfo.Send(ctx)
	}
}
