// Package bilibiliparse B站分享解析
package bilibiliparse

import (
	"fmt"
	"strings"

	logger "github.com/sirupsen/logrus"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

// VideoInfo 视频信息结构体
type VideoInfo struct {
	Title    string
	CoverURL string
	Date     string
	Owner    string
	Like     int
	View     int
	Favorite int
	Coin     int
	Share    int
	Desc     string
	URL      string
	BVID     string
}

// send 发送视频信息
func (videoInfo VideoInfo) send(ctx *zero.Ctx) {
	var messageBuilder strings.Builder
	messageBuilder.WriteString(videoInfo.Title)
	messageBuilder.WriteString("\nUP主: ")
	messageBuilder.WriteString(videoInfo.Owner)
	messageBuilder.WriteString("\n投稿日期: ")
	messageBuilder.WriteString(videoInfo.Date)
	messageBuilder.WriteString("\n点赞: ")
	messageBuilder.WriteString(formatDigit(videoInfo.Like))
	messageBuilder.WriteString("  播放: ")
	messageBuilder.WriteString(formatDigit(videoInfo.View))
	messageBuilder.WriteString("  收藏: ")
	messageBuilder.WriteString(formatDigit(videoInfo.Favorite))
	messageBuilder.WriteString("  硬币: ")
	messageBuilder.WriteString(formatDigit(videoInfo.Coin))
	messageBuilder.WriteString("  分享: ")
	messageBuilder.WriteString(formatDigit(videoInfo.Share))
	messageBuilder.WriteRune('\n')
	messageBuilder.WriteString(videoInfo.Desc)
	messageBuilder.WriteRune('\n')
	messageBuilder.WriteString(videoInfo.URL)

	ctx.SendChain(
		message.Image(videoInfo.CoverURL),
		message.Text(messageBuilder.String()),
	)
	logger.Infoln("已发送视频信息:", videoInfo.BVID)
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
