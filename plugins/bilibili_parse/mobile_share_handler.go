/*
@Title        mobile_share_handler.go
@Description  移动端分享信息处理
@Author       WhitePaper233 2022.7.13
@Update       WhitePaper233 2022.7.13
*/
package bilibili_parse

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	logger "github.com/sirupsen/logrus"
	zero "github.com/wdvxdr1123/ZeroBot"
)

func HandleMobileShare(ctx *zero.Ctx) {
	var appInfo MiniAppInfo
	// 判断是否是移动端分享信息
	if appInfo.IsBilibiliShare(ctx.Event.RawMessage) {
		logger.Debugln("匹配移动端分享信息成功,MessageId:", ctx.Event.MessageID)

		// 获取重定向链接
		redictedLink, err := appInfo.GetRedictLink()
		if err != nil {
			logger.Errorln("获取视频链接失败:", err)
			return
		}

		// 匹配结果
		results := CompiledVIDRegex.FindStringSubmatch(redictedLink)
		vid := results[0]
		logger.Debugln("获取视频ID成功:", vid)

		// 获取视频信息
		videoInfo, err := GetVideoInfo(vid)
		if err != nil {
			logger.Errorln("获取视频信息失败:", err)
			return
		}
		videoInfo.Send(ctx)
	}
}

type MiniAppInfo struct {
	App   string `json:"app"`
	Extra struct {
		AppID int `json:"appid"`
	} `json:"extra"`
	Meta struct {
		Detail_1 struct {
			Qqdocurl string `json:"qqdocurl"`
		} `json:"detail_1"`
	} `json:"meta"`
}

func (appInfo *MiniAppInfo) IsBilibiliShare(rawMessage string) bool {
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
	} else {
		return false
	}
}

func (appInfo *MiniAppInfo) GetRedictLink() (string, error) {
	// 构造请求
	request, err := http.NewRequest("GET", appInfo.Meta.Detail_1.Qqdocurl, nil)
	if err != nil {
		return "", err
	}

	//构造一个禁止重定向的client
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
