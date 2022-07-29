// Package control 插件控制模块
package control

import (
	"time"

	logger "github.com/sirupsen/logrus"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"gorm.io/gorm"

	"github.com/DaydreamCafe/Cocoa/V2/src/config"
	"github.com/DaydreamCafe/Cocoa/V2/src/conn"
	"github.com/DaydreamCafe/Cocoa/V2/src/model"
)

// CheckPremissionHandler 包装handler, 使其拥有全局用户权限鉴权
func CheckPremissionHandler(handler zero.Handler, minLevel int64, echoLevel EchoLevel) zero.Handler {
	return func(ctx *zero.Ctx) {
		// 读取配置文件
		cfg := config.Config{}
		cfg.Load()

		// 连接到数据库
		db, err := conn.GetDB()
		if err != nil {
			logger.Error("用户鉴权失败:", err)
			return
		}

		sqlDB, err := db.DB()
		if err != nil {
			logger.Error("用户鉴权失败:", err)
			return
		}
		defer sqlDB.Close()

		// 查询用户封禁
		// 数据表判空
		var count int64
		err = db.Model(&model.BanedUserModel{}).Count(&count).Error
		if err != nil {
			logger.Error("用户鉴权失败:", err)
			return
		}

		// 数据表不为空
		if count > 0 {
			// 查询用户是否被封禁
			var banedUser model.BanedUserModel
			err := db.Where("qid = ?", ctx.Event.UserID).First(&banedUser).Error
			if err != nil && err != gorm.ErrRecordNotFound {
				logger.Error("用户鉴权失败:", err)
				return
			}

			// 用户是否被封禁
			if err == nil {
				// 判断用户是否在封禁时间内
				if banedUser.UnbanTimeStamp > time.Now().Unix() {
					if echoLevel == 2 || echoLevel == 3 {
						ctx.SendChain(message.Text("您没有权限执行此命令"))
					}
					return
				}

				// 删除用户封禁记录
				err = db.Delete(&banedUser).Error
				if err != nil {
					logger.Error("用户鉴权失败:", err)
					return
				}
			}
		}

		// 查询用户权限
		// 数据表判空
		err = db.Model(&model.UserPremissionModel{}).Count(&count).Error
		if err != nil {
			logger.Error("用户鉴权失败:", err)
			return
		}

		// 数据表无记录
		if count == 0 {
			// 直接与默认等级比较
			if cfg.DefaultLevel >= minLevel {
				handler(ctx)
				return
			}
			if echoLevel == 2 || echoLevel == 3 {
				ctx.SendChain(message.Text("您没有权限执行此命令"))
			}
			return
		}

		// 查询用户
		var user model.UserPremissionModel
		err = db.Where("qid = ?", ctx.Event.UserID).First(&user).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			logger.Error("用户鉴权失败:", err)
			return
		}

		// 当无用户记录时
		if err == gorm.ErrRecordNotFound {
			if cfg.DefaultLevel >= minLevel {
				handler(ctx)
				return
			}
			if echoLevel == 2 || echoLevel == 3 {
				ctx.SendChain(message.Text("您没有权限执行此命令"))
			}
			return
		}

		// 当有用户记录时
		// 判断是否为SU
		if user.IfSU {
			handler(ctx)
			return
		}

		// 用户等级比较
		if user.Level >= minLevel {
			handler(ctx)
			return
		}

		// 不满足SU和等级要求
		if echoLevel == 2 || echoLevel == 3 {
			ctx.SendChain(message.Text("您没有权限执行此命令"))
		}
	}
}

// CheckPremission 用户权限鉴权
func CheckPremission(QID int64, minLevel int64) bool {
	// 读取配置文件
	cfg := config.Config{}
	cfg.Load()

	// 连接到数据库
	db, err := conn.GetDB()
	if err != nil {
		logger.Error("用户鉴权失败:", err)
		return false
	}

	sqlDB, err := db.DB()
	if err != nil {
		logger.Error("用户鉴权失败:", err)
		return false
	}
	defer sqlDB.Close()

	// 查询用户封禁
	// 数据表判空
	var count int64
	err = db.Model(&model.BanedUserModel{}).Count(&count).Error
	if err != nil {
		logger.Error("用户鉴权失败:", err)
		return false
	}

	// 数据表不为空
	if count > 0 {
		// 查询用户是否被封禁
		var banedUser model.BanedUserModel
		err := db.Where("qid = ?", QID).First(&banedUser).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			logger.Error("用户鉴权失败:", err)
			return false
		}

		// 用户是否被封禁
		if err == nil {
			// 判断用户是否在封禁时间内
			if banedUser.UnbanTimeStamp > time.Now().Unix() {
				return false
			}

			// 删除用户封禁记录
			err = db.Delete(&banedUser).Error
			if err != nil {
				logger.Error("用户鉴权失败:", err)
				return false
			}
		}
	}

	// 查询用户权限
	// 数据表判空
	err = db.Model(&model.UserPremissionModel{}).Count(&count).Error
	if err != nil {
		logger.Error("用户鉴权失败:", err)
		return false
	}

	// 数据表无记录
	if count == 0 {
		// 直接与默认等级比较
		return cfg.DefaultLevel >= minLevel
	}

	// 查询用户
	var user model.UserPremissionModel
	err = db.Where("qid = ?", QID).First(&user).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		logger.Error("用户鉴权失败:", err)
		return false
	}

	// 当无用户记录时
	if err == gorm.ErrRecordNotFound {
		return cfg.DefaultLevel >= minLevel
	}

	// 当有用户记录时
	// 判断是否为SU
	if user.IfSU {
		return true
	}

	// 用户等级比较
	if user.Level >= minLevel {
		return true
	}

	// 不满足SU和等级要求
	return false
}
