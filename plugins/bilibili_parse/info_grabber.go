// Package bilibiliparse B站分享解析
package bilibiliparse

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	logger "github.com/sirupsen/logrus"
)

// API B站视频信息API
const API = "http://api.bilibili.com/x/web-interface/view"

// APIStruct B站API结构体
type APIStruct struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	TTL     int    `json:"ttl"`
	Data    struct {
		Bvid     string `json:"bvid"`
		Avid     int    `json:"aid"`
		Title    string `json:"title"`
		CoverURL string `json:"pic"`
		Desc     string `json:"desc"`
		Stat     struct {
			View     int `json:"view"`
			Like     int `json:"like"`
			Favorite int `json:"favorite"`
			Coin     int `json:"coin"`
			Share    int `json:"share"`
		} `json:"stat"`
		Owner struct {
			Name string `json:"name"`
		} `json:"owner"`
	} `json:"data"`
}

// GetVideoInfo 获取视频信息
func GetVideoInfo(vid string) (VideoInfo, error) {
	// 获取视频信息
	APIInfo, err := GetVideoInfoByAPI(vid)
	logger.Debugln("正在获取视频信息:", vid)
	if err != nil {
		return VideoInfo{}, err
	}
	if APIInfo.Code != 0 {
		return VideoInfo{}, errors.New(APIInfo.Message)
	}
	logger.Debugln("获取视频信息成功:", vid)

	var videoInfo VideoInfo
	videoInfo.Title = APIInfo.Data.Title
	videoInfo.Owner = APIInfo.Data.Owner.Name
	videoInfo.CoverURL = APIInfo.Data.CoverURL
	videoInfo.Like = APIInfo.Data.Stat.Like
	videoInfo.View = APIInfo.Data.Stat.View
	videoInfo.Favorite = APIInfo.Data.Stat.Favorite
	videoInfo.Coin = APIInfo.Data.Stat.Coin
	videoInfo.Share = APIInfo.Data.Stat.Share
	videoInfo.Desc = APIInfo.Data.Desc
	videoInfo.URL = fmt.Sprintf("https://www.bilibili.com/video/%s", APIInfo.Data.Bvid)
	videoInfo.BVID = APIInfo.Data.Bvid

	// 返回视频信息
	return videoInfo, nil
}

// GetVideoInfoByAPI 通过API获取视频信息
func GetVideoInfoByAPI(vid string) (APIStruct, error) {
	if strings.HasPrefix(vid, "av") || strings.HasPrefix(vid, "AV") {
		videoInfo, err := GetVideoInfoByAVID(strings.TrimPrefix(vid, "av"))
		return videoInfo, err
	} else if strings.HasPrefix(vid, "bv") || strings.HasPrefix(vid, "BV") {
		videoInfo, err := GetVideoInfoByBVID(vid)
		return videoInfo, err
	} else {
		return APIStruct{}, errors.New("invalid video id")
	}
}

// GetVideoInfoByAVID 通过av号获取视频信息
func GetVideoInfoByAVID(avid string) (APIStruct, error) {
	var api APIStruct
	// 获取视频信息
	request, err := http.NewRequest("GET", API+"?aid="+avid, nil)
	if err != nil {
		return api, err
	}
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return api, err
	}
	defer response.Body.Close()

	err = json.NewDecoder(response.Body).Decode(&api)
	if err != nil {
		return api, err
	}
	return api, nil
}

// GetVideoInfoByBVID 通过BV号获取视频信息
func GetVideoInfoByBVID(bvid string) (APIStruct, error) {
	var api APIStruct
	// 获取视频信息
	request, err := http.NewRequest("GET", API+"?bvid="+bvid, nil)
	if err != nil {
		return api, err
	}
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return api, err
	}
	defer response.Body.Close()

	err = json.NewDecoder(response.Body).Decode(&api)
	if err != nil {
		return api, err
	}
	return api, nil
}
