/*
格式化信息相关
*/
package bili_info_disp

import (
	"fmt"

	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

// 视频信息结构体
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
}

func (videoInfo VideoInfo) Send(ctx *zero.Ctx) {
	ctx.SendChain(
		message.Image(videoInfo.CoverURL),
		message.Text(videoInfo.Title+"\n"),
		message.Text(videoInfo.Owner+"\n"),
		message.Text(
			fmt.Sprint(
				fmt.Sprint("点赞: ", fmt.Sprintf("%d", videoInfo.Like), "  "),
				fmt.Sprint("播放: ", fmt.Sprintf("%d", videoInfo.View), "  "),
				fmt.Sprint("收藏: ", fmt.Sprintf("%d", videoInfo.Favorite), "  "),
				fmt.Sprint("硬币: ", fmt.Sprintf("%d", videoInfo.Coin), "  "),
				fmt.Sprint("分享: ", fmt.Sprintf("%d", videoInfo.Share), "\n"),
			),
		),
		message.Text(videoInfo.Desc+"\n"),
		message.Text(videoInfo.URL),
	) 
}