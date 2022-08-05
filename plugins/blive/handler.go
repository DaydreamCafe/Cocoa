// Package blive B站直播搜索
package blive

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
	BiliSearchAPI = "https://api.bilibili.com/x/web-interface/search/type?page=1&search_type=live&keyword=%s"
	// BiliLiveURL 直播间地址
	BiliLiveURL = "https://live.bilibili.com/"
)

type searchResp struct {
	Code int64 `json:"code"`
	Data struct {
		Result struct {
			LiveRoom []struct {
				RoomCover string `json:"cover"`
				LiveTime  string `json:"live_time"`
				Title     string `json:"title"`
				UserCover string `json:"user_cover"`
				Watched   struct {
					Num int64 `json:"num"`
				} `json:"watched_show"`
			} `json:"live_room"`
			LiveUser []struct {
				IsLive   bool   `json:"is_live"`
				UserName string `json:"uname"`
				LiveTime string `json:"live_time"`
				RoomID   int64  `json:"roomid"`
			} `json:"live_user"`
		} `json:"result"`
	} `json:"data"`
}

func handleBlive(ctx *zero.Ctx) {
	// 解析命令参数
	rawCmd := compiledCommandRegex.FindAllStringSubmatch(ctx.MessageString(), -1)
	targetKeyword := string([]byte(rawCmd[0][0])[6:]) // 截去命令前缀"bililive "
	// 请求API
	response, err := http.Get(fmt.Sprintf(BiliSearchAPI, targetKeyword))
	if err != nil {
		logger.Errorln("查询直播用户错误:", err)
		ctx.SendChain(message.Text("查询直播用户错误"))
		return
	}

	// 解析json
	var userResp searchResp
	err = json.NewDecoder(response.Body).Decode(&userResp)
	if err != nil {
		logger.Errorln("解析json错误:", err)
		ctx.SendChain(message.Text("查询直播用户错误"))
		return
	}

	defer response.Body.Close()

	if len(userResp.Data.Result.LiveUser) == 0 {
		logger.Debugln("未查询到用户:", targetKeyword)
		ctx.SendChain(message.Text("未查询到用户:", targetKeyword))
		return
	}

	// 请求API
	response, err = http.Get(fmt.Sprintf(BiliSearchAPI, strconv.FormatUint(uint64(userResp.Data.Result.LiveUser[0].RoomID), 10)))
	if err != nil {
		logger.Errorln("查询直播间错误:", err)
		ctx.SendChain(message.Text("查询直播间错误"))
		return
	}

	// 解析json
	var roomResp searchResp
	err = json.NewDecoder(response.Body).Decode(&roomResp)
	if err != nil {
		logger.Errorln("解析json错误:", err)
		ctx.SendChain(message.Text("查询直播间错误"))
		return
	}

	defer response.Body.Close()

	if len(roomResp.Data.Result.LiveRoom) == 0 {
		logger.Debugln("未查询到用户直播间:", strconv.FormatUint(uint64(userResp.Data.Result.LiveUser[0].RoomID), 10))
	}

	// 格式化回复信息
	var replyBuilder strings.Builder
	replyBuilder.WriteString("主播: ")
	replyBuilder.WriteString(roomResp.Data.Result.LiveUser[0].UserName)
	if userResp.Data.Result.LiveUser[0].IsLive {
		replyBuilder.WriteString("【直播中】\n")
		replyBuilder.WriteString(roomResp.Data.Result.LiveRoom[0].Title)
		replyBuilder.WriteString("\n--------------------\n")
		replyBuilder.WriteString("开播时间: ")
		replyBuilder.WriteString(roomResp.Data.Result.LiveRoom[0].LiveTime)
		replyBuilder.WriteRune('\n')
		replyBuilder.WriteString(strconv.FormatUint(uint64(roomResp.Data.Result.LiveRoom[0].Watched.Num), 10))
		replyBuilder.WriteString("人观看过\n")
	} else {
		replyBuilder.WriteString("【未开播】\n")
	}
	replyBuilder.WriteString(BiliLiveURL)
	replyBuilder.WriteString(strconv.FormatUint(uint64(userResp.Data.Result.LiveUser[0].RoomID), 10))
	reply := replyBuilder.String()
	// 发送信息
	if userResp.Data.Result.LiveUser[0].IsLive {
		ctx.SendChain(
			message.Image("https:"+roomResp.Data.Result.LiveRoom[0].UserCover),
			message.Text("\n"+reply),
		)
	} else {
		ctx.SendChain(message.Text(reply))
	}
}
