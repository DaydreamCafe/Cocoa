// Package lolicon lolicon图片
package lolicon

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	logger "github.com/sirupsen/logrus"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

const (
	// UA User-Agent
	UA = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.132 Safari/537.36"
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
	//从API获取图片地址
	request, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		logger.Errorln("创建请求失败:", err)
		return
	}
	request.Header.Set("User-Agent", UA)
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		logger.Errorln("请求失败:", err)
		return
	}
	defer response.Body.Close()

	// 将请求结果JSON解析为APIResp结构体
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		logger.Errorln("读取响应失败:", err)
		return
	}

	var resp APIResp
	err = json.Unmarshal(responseBody, &resp)
	if err != nil {
		logger.Errorln("JSON解析失败:", err)
		return
	}
	if resp.Error != "" {
		logger.Errorln("API请求失败:", err)
		return
	}

	// 发送图片
	msg := message.Image(resp.Data[0].URLs.Original)
	rsp := ctx.CallAction("send_group_msg", zero.Params{
		"group_id": ctx.Event.GroupID,
		"message": msg,
	}).Data.Get("message_id")
	
	if rsp.Exists() {
		logger.Infof("发送群消息(%v): [CQ:image,file=%v] (id=%v)", ctx.Event.GroupID, resp.Data[0].URLs.Original, rsp.Int())

		// 撤回图片
		time.Sleep(deleteDelay * time.Second)
		ctx.DeleteMessage(message.NewMessageIDFromInteger(rsp.Int()))
		return
	}

	ctx.SendChain(message.Image("https://pic.imgdb.cn/item/62e3afbdf54cd3f937d4b1af.jpg"))
}
