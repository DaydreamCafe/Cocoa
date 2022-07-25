// Package usermanager 用户管理相关代码
package usermanager

import (
	"errors"
	"flag"
	"fmt"
	"time"

	logger "github.com/sirupsen/logrus"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/extension/shell"
	"github.com/wdvxdr1123/ZeroBot/message"
	"gorm.io/gorm"

	"github.com/DaydreamCafe/Cocoa/V2/src/config"
	"github.com/DaydreamCafe/Cocoa/V2/src/conn"
	"github.com/DaydreamCafe/Cocoa/V2/src/model"
	"github.com/DaydreamCafe/Cocoa/V2/utils/control"
)

func handleUser(ctx *zero.Ctx) {
	cfg := config.Config{}
	cfg.Load()

	fset := flag.FlagSet{}

	// user
	var user int64
	fset.Int64Var(&user, "u", -1, "选择受操作用户")
	fset.Int64Var(&user, "user", -1, "选择受操作用户")

	// set
	var set int64
	fset.Int64Var(&set, "s", -1, "设置用户等级")
	fset.Int64Var(&set, "set", -1, "设置用户等级")

	// ban
	var ban int64
	fset.Int64Var(&ban, "b", -1, "封禁用户")
	fset.Int64Var(&ban, "ban", -1, "封禁用户")

	// pardon
	var pardon bool
	fset.BoolVar(&pardon, "p", false, "解封用户")
	fset.BoolVar(&pardon, "pardon", false, "解封用户")

	// reset
	var reset bool
	fset.BoolVar(&reset, "r", false, "重置用户")
	fset.BoolVar(&reset, "reset", false, "重置用户")

	// help
	var help bool
	fset.BoolVar(&help, "h", false, "帮助")
	fset.BoolVar(&help, "help", false, "帮助")

	args := shell.Parse(ctx.State["args"].(string))
	fset.Parse(args)

	// 执行命令
	// 处理help命令
	if help {
		if control.CheckPremission(ctx.Event.UserID, 5) {
			ctx.SendChain(message.Text(usage))
			return
		}
		ctx.SendChain(message.Text("你没有权限执行这条指令"))
		return
	}

	// 处理user命令
	if user == -1 {
		if control.CheckPremission(ctx.Event.UserID, 9) {
			ctx.SendChain(message.Text("请指定用户"))
			return
		}
		ctx.SendChain(message.Text("你没有权限执行这条指令"))
		return
	}

	// 处理set命令
	if set > 0 {
		if control.CheckPremission(ctx.Event.UserID, 9) {
			err := setPermissionLevel(user, set)
			if err != nil {
				ctx.SendChain(message.Text(err.Error()))
				return
			}
			ctx.SendChain(message.Text(fmt.Sprintf("用户%d的权限等级已设置为%d", user, set)))
			return
		}
		ctx.SendChain(message.Text("你没有权限执行这条指令"))
		return
	}


	// 处理unban命令
	if pardon {
		if control.CheckPremission(ctx.Event.UserID, 9) {
			err := pardonUser(user)
			if err != nil {
				ctx.SendChain(message.Text(err.Error()))
				return
			}
			ctx.SendChain(message.Text(fmt.Sprintf("用户%d已解封", user)))
			return
		}
		ctx.SendChain(message.Text("你没有权限执行这条指令"))
		return
	}

	// 处理reset命令
	if reset {
		if control.CheckPremission(ctx.Event.UserID, 9) {
			err := resetUser(user)
			if err != nil {
				ctx.SendChain(message.Text(err.Error()))
				return
			}
			ctx.SendChain(message.Text(fmt.Sprintf("用户%d的权限等级已重置为%d", user, cfg.DefaultLevel)))
			return
		}
		ctx.SendChain(message.Text("你没有权限执行这条指令"))
		return
	}

	// 处理ban命令
	if control.CheckPremission(ctx.Event.UserID, 9) {
		err := banUser(user, ban)
		if err != nil {
			ctx.SendChain(message.Text(err.Error()))
			return
		}
		ctx.SendChain(message.Text(fmt.Sprintf("用户%d已被封禁", user)))
	}
	ctx.SendChain(message.Text("你没有权限执行这条指令"))
}

// setPermissionLevel 设置用户等级
func setPermissionLevel(QID int64, targetLevel int64) error {
	// 连接到数据库
	db, err := conn.GetDB()
	if err != nil {
		logger.Errorln("更新用户权限失败: ", err)
		return errors.New("指令执行失败: 无法连接到数据库")
	}

	sqlDB, err := db.DB()
	if err != nil {
		logger.Errorln("更新用户权限失败: ", err)
		return errors.New("指令执行失败: 无法连接到数据库")
	}
	defer sqlDB.Close()

	// 查询表中是否有记录
	var count int64
	err = db.Model(&model.UserPremissionModel{}).Count(&count).Error
	if err != nil {
		logger.Errorln("更新用户权限失败: ", err)
		return errors.New("指令执行失败: 查询数据库失败")
	}

	// 如果没有记录, 则创建一条记录
	if count == 0 {
		err = db.Create(&model.UserPremissionModel{
			QID:   QID,
			Level: targetLevel,
			IfSU:  false,
		}).Error
		if err != nil {
			logger.Errorln("更新用户权限失败: ", err)
			return errors.New("指令执行失败: 更新数据库失败")
		}
		return nil
	}

	// 如果有记录, 则查询是否有该用户的记录
	var user model.UserPremissionModel
	err = db.Where("qid = ?", QID).First(&user).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		logger.Errorln("更新用户权限失败: ", err)
		return errors.New("指令执行失败: 查询数据库失败")
	}

	// 如果记录不存在, 则创建一条记录
	if err == gorm.ErrRecordNotFound {
		err = db.Create(&model.UserPremissionModel{
			QID:   QID,
			Level: targetLevel,
			IfSU:  false,
		}).Error
		if err != nil {
			logger.Errorln("更新用户权限失败: ", err)
			return errors.New("指令执行失败: 更新数据库失败")
		}
		return nil
	}

	// 如果记录存在, 则更新记录
	err = db.Model(&user).Update("level", targetLevel).Error
	if err != nil {
		logger.Errorln("更新用户权限失败: ", err)
		return errors.New("指令执行失败: 更新数据库失败")
	}
	return nil
}

// banUser 封禁用户
func banUser(QID int64, targetTime int64) error {
	// 连接到数据库
	db, err := conn.GetDB()
	if err != nil {
		logger.Errorln("封禁用户失败: ", err)
		return errors.New("指令执行失败: 无法连接到数据库")
	}

	sqlDB, err := db.DB()
	if err != nil {
		logger.Errorln("封禁用户失败: ", err)
		return errors.New("指令执行失败: 无法连接到数据库")
	}
	defer sqlDB.Close()

	// 判断目标用户是否为SU
	var userPermission model.UserPremissionModel
	// 查询表中是否有记录
	var count int64
	err = db.Model(userPermission).Count(&count).Error
	if err != nil {
		logger.Errorln("解封用户失败: ", err)
		return errors.New("指令执行失败: 查询数据库失败")
	}
	// 当有记录时, 查询是否有该用户的记录
	if count > 0 {
		err = db.Where("qid = ?", QID).First(&userPermission).Error
		if err != nil {
			logger.Errorln("封禁用户失败: ", err)
			return errors.New("指令执行失败: 查询数据库失败")
		}
		if userPermission.IfSU {
			return errors.New("指令执行失败: 您不能封禁超级用户")
		}
	}

	// 查询表中是否有记录
	err = db.Model(&model.UserPremissionModel{}).Count(&count).Error
	if err != nil {
		logger.Errorln("封禁用户失败: ", err)
		return errors.New("指令执行失败: 查询数据库失败")
	}

	// 如果没有记录, 则创建一条记录
	if count == 0 {
		nowTime := time.Now().Unix()
		var unbanTimeStamp int64
		if targetTime == -1 {
			unbanTimeStamp = -1
		} else {
			unbanTimeStamp = nowTime + targetTime
		}
		err = db.Create(&model.BanedUserModel{
			QID:            QID,
			BanTimeStamp:   nowTime,
			UnbanTimeStamp: unbanTimeStamp,
		}).Error
		if err != nil {
			logger.Errorln("封禁用户失败: ", err)
			return errors.New("指令执行失败: 更新数据库失败")
		}
		return nil
	}

	// 如果有记录, 则查询是否有该用户的记录
	var user model.BanedUserModel
	err = db.Where("qid = ?", QID).First(&user).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		logger.Errorln("封禁用户失败: ", err)
		return errors.New("指令执行失败: 查询数据库失败")
	}

	// 如果记录不存在, 则创建一条记录
	if err == gorm.ErrRecordNotFound {
		nowTime := time.Now().Unix()
		var unbanTimeStamp int64
		if targetTime == -1 {
			unbanTimeStamp = -1
		} else {
			unbanTimeStamp = nowTime + targetTime
		}

		err = db.Create(&model.BanedUserModel{
			QID:            QID,
			BanTimeStamp:   nowTime,
			UnbanTimeStamp: unbanTimeStamp,
		}).Error
		if err != nil {
			logger.Errorln("封禁用户失败: ", err)
			return errors.New("指令执行失败: 更新数据库失败")
		}
		return nil
	}

	// 如果记录存在, 则更新记录
	var unbanTimeStamp int64
	if targetTime == -1 {
		unbanTimeStamp = -1
	} else {
		unbanTimeStamp = time.Now().Unix() + targetTime
	}
	err = db.Model(&user).Update("unban_time_stamp", unbanTimeStamp).Error
	if err != nil {
		logger.Errorln("封禁用户失败: ", err)
		return errors.New("指令执行失败: 更新数据库失败")
	}
	return nil
}

// pardonUser 解封用户
func pardonUser(QID int64) error {
	// 连接到数据库
	db, err := conn.GetDB()
	if err != nil {
		logger.Errorln("解封用户失败: ", err)
		return errors.New("指令执行失败: 无法连接到数据库")
	}

	sqlDB, err := db.DB()
	if err != nil {
		logger.Errorln("解封用户失败: ", err)
		return errors.New("指令执行失败: 无法连接到数据库")
	}
	defer sqlDB.Close()

	// 判断目标用户是否为SU
	var userPermission model.UserPremissionModel
	// 查询表中是否有记录
	var count int64
	err = db.Model(userPermission).Count(&count).Error
	if err != nil {
		logger.Errorln("解封用户失败: ", err)
		return errors.New("指令执行失败: 查询数据库失败")
	}
	// 当有记录时, 查询是否有该用户的记录
	if count > 0 {
		err = db.Where("qid = ?", QID).First(&userPermission).Error
		if err != nil {
			logger.Errorln("解封用户失败: ", err)
			return errors.New("指令执行失败: 查询数据库失败")
		}
		if userPermission.IfSU {
			return errors.New("指令执行失败: 您不能解封超级用户")
		}
	}

	// 查询表中是否有记录
	err = db.Model(&model.BanedUserModel{}).Count(&count).Error
	if err != nil {
		logger.Errorln("解封用户失败: ", err)
		return errors.New("指令执行失败: 查询数据库失败")
	}

	// 如果没有记录, 则返回用户未被封禁
	if count == 0 {
		return errors.New("指令执行失败: 用户未被封禁")
	}

	// 如果有记录, 则查询是否有该用户的记录
	var user model.BanedUserModel
	err = db.Where("qid = ?", QID).First(&user).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		logger.Errorln("解封用户失败: ", err)
		return errors.New("指令执行失败: 查询数据库失败")
	}

	// 如果记录不存在, 则返回用户未被封禁
	if err == gorm.ErrRecordNotFound {
		return errors.New("指令执行失败: 用户未被封禁")
	}

	// 如果记录存在, 则删除记录
	err = db.Delete(&user).Error
	if err != nil {
		logger.Errorln("解封用户失败: ", err)
		return errors.New("指令执行失败: 更新数据库失败")
	}
	return nil
}

// resetUser 重置用户等级和封禁时间
func resetUser(QID int64) error {
	// 读取配置文件
	cfg := config.Config{}
	cfg.Load()

	// 连接到数据库
	db, err := conn.GetDB()
	if err != nil {
		logger.Errorln("重置用户等级和封禁时间失败: ", err)
		return errors.New("指令执行失败: 无法连接到数据库")
	}

	sqlDB, err := db.DB()
	if err != nil {
		logger.Errorln("重置用户等级和封禁时间失败: ", err)
		return errors.New("指令执行失败: 无法连接到数据库")
	}
	defer sqlDB.Close()

	// 重置用户等级
	var userPermission model.UserPremissionModel
	// 查询表中是否有记录
	var count int64
	err = db.Model(userPermission).Count(&count).Error
	if err != nil {
		logger.Errorln("重置用户等级和封禁时间失败: ", err)
		return errors.New("指令执行失败: 查询数据库失败")
	}
	// 当有记录时, 查询是否有该用户的记录
	if count > 0 {
		err = db.Where("qid = ?", QID).First(&userPermission).Error
		if err != nil {
			logger.Errorln("重置用户等级和封禁时间失败: ", err)
			return errors.New("指令执行失败: 查询数据库失败")
		}
		// 判断目标用户是否为SU
		if userPermission.IfSU {
			return errors.New("指令执行失败: 您不能重置超级用户")
		}

		// 若目标用户为普通用户, 则重置用户等级
		err = db.Model(&userPermission).Update("level", cfg.DefaultLevel).Error
		if err != nil {
			logger.Errorln("重置用户等级和封禁时间失败: ", err)
			return errors.New("指令执行失败: 更新数据库失败")
		}
		return nil
	}

	// 重置用户封禁时间
	var user model.BanedUserModel
	// 查询表中是否有记录
	err = db.Model(user).Count(&count).Error
	if err != nil {
		logger.Errorln("重置用户等级和封禁时间失败: ", err)
		return errors.New("指令执行失败: 查询数据库失败")
	}

	// 当有记录时, 查询是否有该用户的记录
	// 因为上面已经确认过不是SU用户, 所以不需要再判断
	if count > 0 {
		err = db.Where("qid = ?", QID).First(&user).Error
		if err != nil {
			logger.Errorln("重置用户等级和封禁时间失败: ", err)
			return errors.New("指令执行失败: 查询数据库失败")
		}

		// 重置用户封禁时间
		err = db.Delete(&user).Error
		if err != nil {
			logger.Errorln("重置用户等级和封禁时间失败: ", err)
			return errors.New("指令执行失败: 更新数据库失败")
		}
		return nil
	}

	// 如果表中没有记录, 那么就什么都不做就等价于重置掉啦（
	return nil
}
