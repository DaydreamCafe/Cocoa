/*
移动端分享相关
*/
package bili_info_disp

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	logger "github.com/sirupsen/logrus"

	_ "github.com/DaydreamCafe/Cocoa/V2/src/logger"
)

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
	var redictedLink string

	// 设置请求头
	request.Header.Set("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 11_0 like Mac OS X) AppleWebKit/604.1.38 (KHTML, like Gecko) Version/11.0 Mobile/15A372 Safari/604.1")

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
	logger.Debugln("捕捉到视频:", string(body))
	// 抽取重定向链接
	redictPage := string(body)
	redictedLink = strings.TrimPrefix(redictPage, "<a href=\"")
	redictedLink = strings.TrimSuffix(redictedLink, "\">Found</a>.")

	return redictedLink, nil
}
