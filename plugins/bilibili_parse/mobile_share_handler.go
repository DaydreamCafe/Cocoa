// Package bilibiliparse B站分享解析
package bilibiliparse

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	logger "github.com/sirupsen/logrus"
	zero "github.com/wdvxdr1123/ZeroBot"
)

// handleMobileShare 移动端分享handler
func handleMobileShare(ctx *zero.Ctx) {
	var appInfo MiniAppInfo
	// 判断是否是移动端分享信息
	if appInfo.isBilibiliShare(ctx.Event.RawMessage) {
		logger.Debugln("匹配移动端分享信息成功,MessageId:", ctx.Event.MessageID)

		// 获取重定向链接
		redictedLink, err := appInfo.getRedictLink()
		if err != nil {
			logger.Errorln("获取视频链接失败:", err)
			return
		}

		// 匹配结果
		results := compiledVIDRegex.FindStringSubmatch(redictedLink)
		vid := results[0]
		logger.Debugln("获取视频ID成功:", vid)

		// 获取视频信息
		videoInfo, err := getVideoInfo(vid)
		if err != nil {
			logger.Errorln("获取视频信息失败:", err)
			return
		}
		videoInfo.send(ctx)
	}
}

// MiniAppInfo 小程序卡片信息结构体
type MiniAppInfo struct {
	App   string `json:"app"`
	Extra struct {
		AppID int `json:"appid"`
	} `json:"extra"`
	Meta struct {
		Detail1 struct {
			Qqdocurl string `json:"qqdocurl"`
		} `json:"detail_1"`
	} `json:"meta"`
}

// isBilibiliShare 判断是否是移动端分享信息
func (appInfo *MiniAppInfo) isBilibiliShare(rawMessage string) bool {
	// 判断是否为json类型的消息段
	if !strings.HasPrefix(rawMessage, "[CQ:json,") {
		return false
	}

	// 预处理消息
	rawMessage = strings.TrimPrefix(rawMessage, "[CQ:json,data=")
	rawMessage = strings.TrimSuffix(rawMessage, "]")
	rawMessage = strings.ReplaceAll(rawMessage, "&#44;", ",")
	// 反序列化消息
	err := json.Unmarshal([]byte(rawMessage), appInfo)
	if err != nil {
		return false
	}
	// 判断是否属于移动端分享
	if appInfo.Extra.AppID == 100951776 {
		return true
	}
	return false
}

// getRedictLink 获取重定向链接
func (appInfo *MiniAppInfo) getRedictLink() (string, error) {
	// 构造请求
	request, err := http.NewRequest("GET", appInfo.Meta.Detail1.Qqdocurl, nil)
	if err != nil {
		return "", err
	}

	// 构造一个禁止重定向的client
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) > 0 {
				return http.ErrUseLastResponse
			}
			return nil
		},
	}

	// 发送请求
	response, err := client.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	// 解析响应为字符串
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	// 抽取重定向链接
	redictPage := string(body)

	var redictedLink string
	redictedLink = strings.TrimPrefix(redictPage, "<a href=\"")
	redictedLink = strings.TrimSuffix(redictedLink, "\">Found</a>.")

	return redictedLink, nil
}
