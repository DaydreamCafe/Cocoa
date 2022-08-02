// Package control 插件控制模块
package control

import (
	"fmt"
	"strconv"
	"strings"

	logger "github.com/sirupsen/logrus"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"gorm.io/gorm"

	"github.com/DaydreamCafe/Cocoa/V2/src/config"
	"github.com/DaydreamCafe/Cocoa/V2/src/conn"
	"github.com/DaydreamCafe/Cocoa/V2/src/model"
)

// 通用全局忽略账号prehandler, 用于确认账号是否被设置为忽略
func ignoreUserChecker(ctx *zero.Ctx) bool {
	// 加载配置文件
	cfg := config.Config{}
	cfg.Load()

	// 检查该账号是否被忽略
	flag := true
	for _, userID := range(cfg.IgnoreUsers) {
		if userID == ctx.Event.UserID {
			flag = false
		}
	}

	return flag
}

// 通用的全局插件prehandler, 用于确认插件是否处于可用状态
func pluginChecker(pluginMetadata Metadata, echoLevel EchoLevel) zero.Rule {
	return func(ctx *zero.Ctx) bool {
		// 连接数据库
		db, err := conn.GetDB()
		if err != nil {
			logger.Errorln("获取数据库连接失败:", err)
			if echoLevel == 1 || echoLevel == 3 {
				ctx.SendChain(message.Text("插件功能错误: 数据库错误"))
			}
			return false
		}

		sqlDB, err := db.DB()
		if err != nil {
			logger.Errorln("获取数据库连接失败:", err)
			if echoLevel == 1 || echoLevel == 3 {
				ctx.SendChain(message.Text("插件功能错误: 数据库错误"))
			}
			return false
		}
		defer sqlDB.Close()

		/*
			检测是否被局部禁用
		*/
		// 查询表中是否有记录
		var localPlugin model.LocalPluginManagement
		var count int64
		err = db.Model(&localPlugin).Count(&count).Error
		if err != nil {
			logger.Errorln("查询数据库失败:", err)
			if echoLevel == 1 || echoLevel == 3 {
				ctx.SendChain(message.Text("插件功能错误: 数据库错误"))
			}
			return false
		}

		// 当表中有记录时, 查询该群是否在记录中
		if count > 0 {
			// 查询插件是否被局部禁用
			err := db.Where("name = ?", pluginMetadata.Name).First(&localPlugin).Error
			if err != nil && err != gorm.ErrRecordNotFound {
				logger.Errorln("查询插件是否被局部禁用失败:", err)
				if echoLevel == 1 || echoLevel == 3 {
					ctx.SendChain(message.Text("插件功能错误: 数据库错误"))
				}
				return false
			}

			// 当插件有记录时, 查询该群是否在记录中
			if err != gorm.ErrRecordNotFound {
				// 当插件被局部禁用时, 返回false
				banedGroups := strings.Split(localPlugin.BanedGroupID, "|")
				currentGroupID := strconv.FormatInt(ctx.Event.GroupID, 10)
				for _, groupID := range banedGroups {
					if groupID == currentGroupID {
						if echoLevel == 1 || echoLevel == 3 {
							ctx.SendChain(message.Text(fmt.Sprintf("%s插件已在该群被禁用", pluginMetadata.Name)))
						}
						return false
					}
				}
			}
		}

		/*
			检测是否被全局禁用
		*/
		var globalPlugin model.GlobalPluginManagement
		err = db.Where("name = ?", pluginMetadata.Name).First(&globalPlugin).Error
		if err != nil {
			logger.Errorln("查询插件是否被全部禁用失败:", err)
			if echoLevel == 1 || echoLevel == 3 {
				ctx.SendChain(message.Text("插件功能错误: 数据库错误"))
			}
			return false
		}

		// 当插件被全部禁用时, 返回false
		if globalPlugin.IsBaned {
			if echoLevel == 1 || echoLevel == 3 {
				ctx.SendChain(message.Text(fmt.Sprintf("%s插件已被全局禁用", pluginMetadata.Name)))
			}
			return false
		}

		// 否则返回true
		return true
	}
}