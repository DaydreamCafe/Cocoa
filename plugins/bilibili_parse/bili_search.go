// Package bilibiliparse B站分享解析
package bilibiliparse

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	logger "github.com/sirupsen/logrus"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

const (
	// BiliSearchAPI API地址
	BiliSearchAPI = "https://api.bilibili.com/x/web-interface/search/type?page=%s&search_type=%s&keyword="
	// Bilibili直播间 直播间地址
	BiliLiveURL = "https://live.bilibili.com/"
)

// BiliLiveSearchAPIResp Bilibili直播搜索返回结构体
type BiliLiveSearchAPIResp struct {
	Code int `json:"code"`
	Data struct {
		Result struct {
			LiveRoom []struct {
				IsOnline bool   `json:"is_live_room_online"`
				RoomID   int    `json:"room_id"`
				Title    string `json:"title"`
				UID      int    `json:"uid"`
				UserName string `json:"uname"`
				CoverURL string `json:"cover"`
				Watched  struct {
					Num int `json:"num"`
				} `json:"watched_show"`
			} `json:"live_room"`
		} `json:"result"`
	} `json:"data"`
}

// LiveSearch 直播间搜索
func handleLiveSearch(ctx *zero.Ctx) {
	// 默认API参数
	LiveSearchAPI := fmt.Sprintf(BiliSearchAPI, "1", "live")
	// 处理搜索参数
	raw_cmd := compiledLiveSearchRegex.FindAllStringSubmatch(ctx.MessageString(), -1)
	cmd_str := string([]byte(raw_cmd[0][0])[8:]) // 截去命令前缀"bililive "
	args := strings.Split(cmd_str, "&")
	var reqURLs = make([]string, len(args))
	for index, arg := range args {
		var reqURLBuilder strings.Builder

		reqURLBuilder.WriteString(LiveSearchAPI)
		reqURLBuilder.WriteString(arg)
		reqURLs[index] = reqURLBuilder.String()
	}

	// 从API获取直播间信息，并解析为结构体
	var respInfos = make([]BiliLiveSearchAPIResp, len(reqURLs))
	for index, reqURL := range reqURLs {
		// 调用API获取信息
		response, err := http.Get(reqURL)
		if err != nil {
			logger.Errorln("请求失败:", err)
			continue
		}
		defer response.Body.Close()

		// 将请求结果JSON解析为BiliLiveAPIResp结构体
		var respInfo BiliLiveSearchAPIResp
		err = json.NewDecoder(response.Body).Decode(&respInfo)
		if err != nil {
			logger.Errorln("JSON解析失败:", err)
			continue
		}

		// 请求失败则跳过
		if respInfo.Code == 0 {
			continue
		}

		respInfos[index] = respInfo
	}

	// 格式化回复字符串及封面图
	var replyStr = make([]string, len(respInfos))
	var coverLst = make([]string, len(respInfos))
	for index, respInfo := range respInfos {
		// 格式化字符串
		var replyBuilder strings.Builder

		replyBuilder.WriteRune('\n')
		replyBuilder.WriteString(respInfo.Data.Result.LiveRoom[0].Title)
		replyBuilder.WriteString("\n主播: ")
		replyBuilder.WriteString(respInfo.Data.Result.LiveRoom[0].UserName)
		replyBuilder.WriteRune('\n')
		replyBuilder.WriteString("--------------------\n")
		replyBuilder.WriteString(strconv.FormatUint(uint64(respInfo.Data.Result.LiveRoom[0].Watched.Num), 10))
		replyBuilder.WriteString("人观看过\n")
		if respInfo.Data.Result.LiveRoom[0].IsOnline {
			replyBuilder.WriteString("【直播中】\n")
		} else {
			replyBuilder.WriteString("【未开播】\n")
		}
		replyBuilder.WriteString(BiliLiveURL)
		replyBuilder.WriteString(formatDigit(respInfo.Data.Result.LiveRoom[0].RoomID))

		replyStr[index] = replyBuilder.String()
		//格式化封面图URL
		coverLst[index] = respInfo.Data.Result.LiveRoom[0].CoverURL
	}

	// 发送信息
	for index, reply := range replyStr {
		// 过滤掉空字符串
		if reply != "" {
			ctx.SendChain(
				message.Image(coverLst[index]),
				message.Text(reply),
			)
		}
	}
}
