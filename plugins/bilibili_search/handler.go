// Package bilibilisearch B站综合搜索
package bilibilisearch

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
	// UA User-Agent
	UA = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.132 Safari/537.36"
	// BiliSearchAPI API地址
	BiliSearchAPI = "https://api.bilibili.com/x/web-interface/search/type?page=%s&search_type=%s&keyword="
	// Bilibili直播间 直播间地址
	BiliLiveURL = "https://live.bilibili.com/"
)

// BiliLiveUserSearchAPIResp Bilibili直播用户搜索返回结构体
type BiliLiveUserSearchAPIResp struct {
	Code int `json:"code"`
	Data struct {
		Result []struct {
			IsLive bool `json:"is_live"`
			RoomID int  `json:"roomid"`
		} `json:"result"`
	} `json:"data"`
}

// BiliLiveSearchAPIResp Bilibili直播搜索返回结构体
type BiliLiveSearchAPIResp struct {
	Code int `json:"code"`
	Data struct {
		Result struct {
			LiveRoom []struct {
				RoomCover string `json:"cover"`
				LiveTime  string `json:"live_time"`
				Title     string `json:"title"`
				UserFace  string `json:"uface"`
				UserCover string `json:"user_cover"`
				UserName  string `json:"uname"`
				Watched   struct {
					Num int `json:"num"`
				} `json:"watchers_show"`
			} `json:"live_room"`
		} `json:"result"`
	} `json:"data"`
}

// LiveSearch 直播间搜索
func handleLiveSearch(ctx *zero.Ctx) {
	// API参数->搜索直播用户
	LiveUserSearchAPI := fmt.Sprintf(BiliSearchAPI, "1", "live_user")
	// 处理搜索参数
	raw_cmd := compiledLiveSearchRegex.FindAllStringSubmatch(ctx.MessageString(), -1)
	cmd_str := string([]byte(raw_cmd[0][0])[9:]) // 截去命令前缀"bililive "
	args := strings.Split(cmd_str, "&amp;")
	logger.Debugln("获取到直播间关键字：", args)
	var reqURLs = make([]string, len(args))
	for index, arg := range args {
		var reqURLBuilder strings.Builder

		reqURLBuilder.WriteString(LiveUserSearchAPI)
		reqURLBuilder.WriteString(arg)
		reqURLs[index] = reqURLBuilder.String()
	}

	// 从API获取直播用户信息，并解析为结构体
	var respInfos_USER = make([]BiliLiveUserSearchAPIResp, len(reqURLs))
	var respInfos_ROOM = make([]BiliLiveSearchAPIResp, len(reqURLs))
	for index, reqURL := range reqURLs {
		// 调用API获取信息
		request, err := http.NewRequest("GET", reqURL, nil)
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

		// 将请求结果JSON解析为BiliLiveAPIResp结构体
		var respInfo_USER BiliLiveUserSearchAPIResp
		err = json.NewDecoder(response.Body).Decode(&respInfo_USER)
		if err != nil {
			logger.Errorln("JSON解析失败:", err)
			continue
		}

		// 请求失败则跳过
		if respInfo_USER.Code != 0 {
			continue
		}

		respInfos_USER[index] = respInfo_USER

		// API参数->搜索直播间
		LiveSearchAPI := fmt.Sprintf(BiliSearchAPI, "1", "live")

		// 依照RoomID重新构建reqURLs
		for index := range args {
			var reqURLBuilder strings.Builder

			reqURLBuilder.WriteString(LiveSearchAPI)
			reqURLBuilder.WriteString(formatDigit(respInfo_USER.Data.Result[0].RoomID))
			reqURLs[index] = reqURLBuilder.String()
		}

		// 从API获取直播间信息，并解析为结构体
		for index, reqURL := range reqURLs {
			// 调用API获取信息
			request, err := http.NewRequest("GET", reqURL, nil)
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

			// 将请求结果JSON解析为BiliLiveAPIResp结构体
			var respInfo_ROOM BiliLiveSearchAPIResp
			err = json.NewDecoder(response.Body).Decode(&respInfo_ROOM)
			if err != nil {
				logger.Errorln("JSON解析失败:", err)
				continue
			}

			// 请求失败则跳过
			if respInfo_ROOM.Code != 0 {
				continue
			}

			respInfos_ROOM[index] = respInfo_ROOM
		}
	}

	// 格式化回复字符串及封面图
	var replyStr = make([]string, len(respInfos_USER))
	var coverLst = make([]string, len(respInfos_ROOM))
	for index, roomInfo := range respInfos_ROOM {
		userInfo := respInfos_USER[index]
		// 格式化字符串
		var replyBuilder strings.Builder

		replyBuilder.WriteString("\n主播: ")
		replyBuilder.WriteString(roomInfo.Data.Result.LiveRoom[0].UserName)
		if userInfo.Data.Result[0].IsLive {
			replyBuilder.WriteString("【直播中】\n")
		} else {
			replyBuilder.WriteString("【未开播】\n")
			continue
		}
		replyBuilder.WriteRune('\n')
		replyBuilder.WriteString(roomInfo.Data.Result.LiveRoom[0].Title)
		replyBuilder.WriteRune('\n')
		replyBuilder.WriteString("--------------------\n")
		replyBuilder.WriteString(strconv.FormatUint(uint64(roomInfo.Data.Result.LiveRoom[0].Watched.Num), 10))
		replyBuilder.WriteString("人观看过\n")
		replyBuilder.WriteString(BiliLiveURL)
		replyBuilder.WriteString(formatDigit(userInfo.Data.Result[0].RoomID))

		replyStr[index] = replyBuilder.String()
		//格式化封面图URL
		coverLst[index] = roomInfo.Data.Result.LiveRoom[0].UserCover
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
