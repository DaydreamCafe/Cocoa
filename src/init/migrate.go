// Package init Bot初始化相关代码
package init

import (
	logger "github.com/sirupsen/logrus"

	"github.com/DaydreamCafe/Cocoa/V2/src/conn"
	"github.com/DaydreamCafe/Cocoa/V2/src/model"
)

func init() {
	db, err := conn.GetDB()
	if err != nil {
		logger.Panicln(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		logger.Panicln(err)
	}
	defer sqlDB.Close()

	db.Exec("DROP TABLE IF EXISTS plugin;")
	// 迁移插件元数据表
	err = db.AutoMigrate(&model.Plugin{})
	if err != nil {
		logger.Panicln("数据库自动迁移失败:", err)
	}

	// 迁移全局插件管理表
	err = db.AutoMigrate(&model.GlobalPluginManagement{})
	if err != nil {
		logger.Panicln("数据库自动迁移失败:", err)
	}

	// 迁移局部插件管理表
	err = db.AutoMigrate(&model.LocalPluginManagement{})
	if err != nil {
		logger.Panicln("数据库自动迁移失败:", err)
	}

	// 迁移用户权限表
	err = db.AutoMigrate(&model.UserPremissionModel{})
	if err != nil {
		logger.Panicln("数据库自动迁移失败:", err)
	}

	// 迁移用户封禁表
	err = db.AutoMigrate(&model.BanedUserModel{})
	if err != nil {
		logger.Panicln("数据库自动迁移失败:", err)
	}
}