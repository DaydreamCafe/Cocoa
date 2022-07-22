// Package bilibiliparse B站分享解析
package bilibiliparse

import (
	"fmt"

	logger "github.com/sirupsen/logrus"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

// VideoInfo 视频信息结构体
type VideoInfo struct {
	Title    string
	CoverURL string
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
	ctx.SendChain(
		message.Image(videoInfo.CoverURL),
		message.Text(videoInfo.Title+"\n"),
		message.Text(videoInfo.Owner+"\n"),
		message.Text(
			fmt.Sprint(
				fmt.Sprint("点赞: ", formatDigit(videoInfo.Like), "  "),
				fmt.Sprint("播放: ", formatDigit(videoInfo.View), "  "),
				fmt.Sprint("收藏: ", formatDigit(videoInfo.Favorite), "  "),
				fmt.Sprint("硬币: ", formatDigit(videoInfo.Coin), "  "),
				fmt.Sprint("分享: ", formatDigit(videoInfo.Share), "\n"),
			),
		),
		message.Text(videoInfo.Desc+"\n"),
		message.Text(videoInfo.URL),
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
