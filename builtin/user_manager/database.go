// Package usermanager 用户管理相关代码
package usermanager

import (
	"time"

	logger "github.com/sirupsen/logrus"

	"github.com/DaydreamCafe/Cocoa/V2/src/conn"
	"github.com/DaydreamCafe/Cocoa/V2/src/model"
	"github.com/DaydreamCafe/Cocoa/V2/src/config"
)


func initDatabase() {
	// 获取数据库连接
	db, err := conn.GetDB()
	if err != nil {
		logger.Panicln("数据库连接失败: ", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		logger.Panicln("数据库连接失败: ", err)
	}
	defer sqlDB.Close()

	/*
		初始化封禁表
	*/
	// 查询封禁表中是否有记录
	var count int64
	err = db.Model(&model.BanedUserModel{}).Count(&count).Error
	if err != nil {
		logger.Panicln("数据库查询失败:", err)
	}

	// 如果有记录, 则删除过时的记录
	if count > 0 {
		// 被封禁用户记录切片
		var banedUsers []model.BanedUserModel
		err = db.Find(&banedUsers).Error
		if err != nil {
			logger.Panicln("数据库查询失败:", err)
		}

		// 遍历被封禁用户记录切片, 删除已到封禁时间的记录
		for _, banedUser := range banedUsers {
			nowTime := time.Now().Unix()
			if banedUser.UnbanTimeStamp != -1 && banedUser.UnbanTimeStamp < nowTime {
				db.Delete(&banedUser)
			}
		}
	}
	
	/*
		初始化用户权限表
	*/
	cfg := config.Config{}
	cfg.Load()

	// 读取用户权限表
	var users []model.UserPremissionModel
	err = db.Find(&users).Error
	if err != nil {
		logger.Panicln("数据库查询失败:", err)
	}

	// 遍历用户权限表, 如果SU不在配置文件中, 则取消SU权限
	for _, user := range users {
		flag := true
		for _, su := range cfg.SuperUsers {
			if user.QID == su {
				flag = false
				break
			}
		}

		// 如果用户不是SU, 则取消SU权限
		if flag {
			err = db.Model(&user).Update("if_su", false).Error
			if err != nil {
				logger.Panicln("数据库更新失败:", err)
			}
		}
	}

	// 遍历配置文件SU列表, 如果在用户权限表中不存在或者没有SU权限, 则添加记录或者更新权限
	for _, su := range cfg.SuperUsers {
		flag := true
		for _, user := range users {
			if user.QID == su {
				flag = false
				err = db.Model(&user).Update("if_su", true).Error
				if err != nil {
					logger.Panicln("数据库更新失败:", err)
				}
				break
			}
		}

		// 如果用户记录不存在, 则添加记录
		if flag {
			user := model.UserPremissionModel{
				QID: su,
				Level: cfg.DefaultLevel,
				IfSU: true,
			}
			err = db.Create(&user).Error
			if err != nil {
				logger.Panicln("数据库添加失败:", err)
			}
		}
	}
}