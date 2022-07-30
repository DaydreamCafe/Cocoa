// Package coser coser图片
package coser

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
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
	apiURL = "http://ovooa.com/API/cosplay/api.php"

	// deleteDelay 撤回图片的延迟
	deleteDelay = 30
)

// coserAPIResp API返回结构体
type coserAPIResp struct {
	Code string `json:"code"`
	Text string `json:"text"`
	Data struct {
		Title string   `json:"Title"`
		Data  []string `json:"data"`
	} `json:"data"`
}

// handleCoser coser命令handler
func handleCoser(ctx *zero.Ctx) {
	// 从API获取图片地址
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

	// 将请求结果JSON解析为CoserApiResp结构体
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		logger.Errorln("读取响应失败:", err)
		return
	}

	var resp coserAPIResp
	err = json.Unmarshal(responseBody, &resp)
	if err != nil {
		logger.Errorln("JSON解析失败:", err)
		return
	}
	if resp.Code != "1" {
		logger.Errorln("API请求失败:", resp.Text)
		return
	}

	// 从所有图片地址中随机选择一个
	var imageURL string
	if len(resp.Data.Data) > 0 {
		// 随机选择一个图片地址
		rand.Seed(time.Now().Unix())
		imageURL = resp.Data.Data[rand.Intn(len(resp.Data.Data))]
	} else {
		logger.Errorln("API返回数据为空")
		return
	}

	// 发送图片
	rsp := ctx.CallAction("send_group_msg", zero.Params{
		"group_id": ctx.Event.GroupID,
		"message": imageURL,
	}).Data.Get("message_id")
	
	if rsp.Exists() {
		logger.Infof("发送群消息(%v): [CQ:image,file=%v] (id=%v)", ctx.Event.GroupID, imageURL, rsp.Int())

		// 撤回图片
		time.Sleep(deleteDelay * time.Second)
		ctx.DeleteMessage(message.NewMessageIDFromInteger(rsp.Int()))
		return
	}

	ctx.SendChain(message.Text("图片发送失败"))
}
