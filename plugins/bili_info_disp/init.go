/*
一个将B站链接/av号/BV号/移动端分享转换为视频信息的插件
*/
package bili_info_disp

import (
	"regexp"

	logger "github.com/sirupsen/logrus"
	zero "github.com/wdvxdr1123/ZeroBot"

	_ "github.com/DaydreamCafe/Cocoa/V2/src/logger"
)

func init() {
	logger.Debugln("BiliInfoDisp 插件加载成功")

	// 编译正则表达式
	compileRegex := regexp.MustCompile(
		`((av|AV)\d+|(bv|BV)1(\d|\w){2}4(\d|\w)1(\d|\w)7(\d|\w){2})`,
	)

	// 匹配av号或者BV号
	zero.OnRegex(`((av|AV)\d+|(bv|BV)1(\d|\w){2}4(\d|\w)1(\d|\w)7(\d|\w){2})`).Handle(
		func(ctx *zero.Ctx) {
			logger.Debugln("匹配AV号或者BV号成功")

			// 匹配结果
			results := compileRegex.FindAllStringSubmatch(ctx.MessageString(), -1)

			// 循环匹配结果
			var VID []string
			for _, result := range results {
				// 获取av号或者BV号
				VID = append(VID, result[0])
			}

			for _, vid := range VID {
				// 获取视频信息
				videoInfo, err := GetVideoInfo(vid)
				if err != nil {
					logger.Errorln("获取视频信息失败: ", err)
					continue
				}
				// 发送视频信息
				videoInfo.Send(ctx)
			}
		},
	)

	// 匹配移动端分享信息
	zero.OnMessage().Handle(
		func(ctx *zero.Ctx) {
			var appInfo MiniAppInfo
			// 判断是否是移动端分享信息
			if appInfo.IsBilibiliShare(ctx.Event.RawMessage) {
				logger.Debugln("匹配移动端分享信息成功")

				// 获取重定向链接
				redictedLink, err := appInfo.GetRedictLink()
				if err != nil {
					logger.Errorln("获取视频链接失败: ", err)
					return
				}
				logger.Debugln("获取视频链接成功: ", redictedLink)

				// 匹配结果
				results := compileRegex.FindStringSubmatch(redictedLink)
				vid := results[0]

				// 获取视频信息
				videoInfo, err := GetVideoInfo(vid)
				if err != nil {
					logger.Errorln("获取视频信息失败: ", err)
					return
				}
				videoInfo.Send(ctx)
			}
		},
	)
}
