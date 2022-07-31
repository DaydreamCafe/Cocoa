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
	var msgBuilder strings.Builder
	msgBuilder.WriteString(videoInfo.Title)
	msgBuilder.WriteString("\nUP主: ")
	msgBuilder.WriteString(videoInfo.Owner)
	msgBuilder.WriteString("\n投稿日期: ")
	msgBuilder.WriteString(videoInfo.Date)
	msgBuilder.WriteString("\n点赞: ")
	msgBuilder.WriteString(formatDigit(videoInfo.Like))
	msgBuilder.WriteString("  播放: ")
	msgBuilder.WriteString(formatDigit(videoInfo.View))
	msgBuilder.WriteString("  收藏: ")
	msgBuilder.WriteString(formatDigit(videoInfo.Favorite))
	msgBuilder.WriteString("  硬币: ")
	msgBuilder.WriteString(formatDigit(videoInfo.Coin))
	msgBuilder.WriteString("  分享: ")
	msgBuilder.WriteString(formatDigit(videoInfo.Share))
	msgBuilder.WriteString("\n简介: ")
	msgBuilder.WriteString(videoInfo.Desc)
	msgBuilder.WriteRune('\n')
	msgBuilder.WriteString(videoInfo.URL)

	ctx.SendChain(
		message.Text(msgBuilder.String()),
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
