// Package lolicon lolicon图片
package lolicon

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"time"

	logger "github.com/sirupsen/logrus"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

const (
	// apiURL API地址
	apiURL = "https://api.lolicon.app/setu/v2"

	// deleteDelay 撤回图片的延迟
	deleteDelay = 30
)

// ImgData 具体图片数据
type ImgData struct {
	PID    int64    `json:"pid"`
	UID    int64    `json:"uid"`
	Title  string   `json:"title"`
	Author string   `json:"author"`
	R18    bool     `json:"r18"`
	Tags   []string `json:"tags"`
	URLs   struct {
		Original string `json:"original"`
	} `json:"urls"`
}

// APIResp LoliconAPI返回结果结构体
type APIResp struct {
	Error string    `json:"error"`
	Data  []ImgData `json:"data"`
}

// handleLoli 涩图命令handler
func handleLoli(ctx *zero.Ctx) {
	// 请求URL
	var reqURL string = apiURL

	// 解析tag
	msg := ctx.ExtractPlainText()
	var tagGroup []string
	if len(msg) > 6 && msg[6:7] == " " {
		tags := ctx.Event.Message.String()[7:]
		tagGroup = strings.Split(tags, "&amp;")

		// 构造请求URL
		var reqURLBuilder strings.Builder
		reqURLBuilder.WriteString(apiURL)
		reqURLBuilder.WriteRune('?')
		for index, orTags := range tagGroup {
			if index != 0 {
				reqURLBuilder.WriteRune('&')
			}
			reqURLBuilder.WriteString("tag=")
			reqURLBuilder.WriteString(url.QueryEscape(orTags))
		}
		reqURL = reqURLBuilder.String()
	}

	//从API获取图片地址
	response, err := http.Get(reqURL)
	if err != nil {
		logger.Errorln("请求失败:", err)
		return
	}
	defer response.Body.Close()

	// 将请求结果JSON解析为APIResp结构体
	var resp APIResp
	err = json.NewDecoder(response.Body).Decode(&resp)
	if err != nil {
		logger.Errorln("JSON解析失败:", err)
		return
	}

	// 发送图片

	// 处理没有该标签涩图的情况
	if len(resp.Data) < 1 {
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("未找到指定标签的涩图"))
		return
	}

	cqcode := message.Image(resp.Data[0].URLs.Original)
	rsp := ctx.CallAction("send_group_msg", zero.Params{
		"group_id": ctx.Event.GroupID,
		"message":  cqcode,
	}).Data.Get("message_id")

	if rsp.Exists() {
		logger.Infof("发送群消息(%v): [CQ:image,file=%v] (id=%v)", ctx.Event.GroupID, resp.Data[0].URLs.Original, rsp.Int())

		// 撤回图片
		time.Sleep(deleteDelay * time.Second)
		ctx.DeleteMessage(message.NewMessageIDFromInteger(rsp.Int()))
		return
	}

	ctx.SendChain(
		message.Reply(ctx.Event.MessageID),
		message.Text("图片发送失败, 你可以自行访问图片链接查看:"),
		message.Text(resp.Data[0].URLs.Original),
	)
}
