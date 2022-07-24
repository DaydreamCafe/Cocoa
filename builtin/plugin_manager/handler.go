// Package pluginmanager 插件管理器
package pluginmanager

import (
	"errors"
	"flag"
	"fmt"
	"strconv"
	"strings"

	logger "github.com/sirupsen/logrus"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/extension/shell"
	"github.com/wdvxdr1123/ZeroBot/message"
	"gorm.io/gorm"

	"github.com/DaydreamCafe/Cocoa/V2/src/conn"
	"github.com/DaydreamCafe/Cocoa/V2/src/model"
	"github.com/DaydreamCafe/Cocoa/V2/utils/control"
)

// handlePlugin 处理插件管理命令
func handlePlugin(ctx *zero.Ctx) {
	fset := flag.FlagSet{}

	// list
	var list bool
	fset.BoolVar(&list, "l", false, "显示所有插件的列表")
	fset.BoolVar(&list, "list", false, "显示所有插件的列表")

	// ban
	var ban string
	fset.StringVar(&ban, "b", "", "禁用插件")
	fset.StringVar(&ban, "ban", "", "禁用插件")

	// unban
	var unban string
	fset.StringVar(&unban, "u", "", "解禁插件")
	fset.StringVar(&unban, "unban", "", "解禁插件")

	// enable
	var enable string
	fset.StringVar(&enable, "e", "", "在本群启用插件")
	fset.StringVar(&enable, "enable", "", "在本群启用插件")

	// disable
	var disable string
	fset.StringVar(&disable, "d", "", "在本群禁用插件")
	fset.StringVar(&disable, "disable", "", "在本群禁用插件")

	// help
	var help bool
	fset.BoolVar(&help, "h", false, "显示该指令的帮助")
	fset.BoolVar(&help, "help", false, "显示该指令的帮助")

	args := shell.Parse(ctx.State["args"].(string))
	fset.Parse(args)

	// 执行指令
	// hasArg 是否提供了参数
	hasArg := false

	// 处理list指令
	if list {
		if control.CheckPremission(ctx.Event.UserID, 5) {
			hasArg = true
			ctx.SendChain(message.Text(getPluginList()))
		} else {
			ctx.SendChain(message.Text("你没有权限执行这条指令"))
			return	
		}
	}

	// 处理ban指令
	if ban != "" {
		if control.CheckPremission(ctx.Event.UserID, 9) {
			hasArg = true
			err := banPlugin(ban)
			if err != nil {
				ctx.SendChain(message.Text(err.Error()))
			} else {
				ctx.SendChain(message.Text((fmt.Sprintf("插件%s全局禁用成功", ban))))
			}
		} else {
			ctx.SendChain(message.Text("你没有权限执行这条指令"))
			return
		}
	}

	// 处理unban指令
	if unban != "" {
		if control.CheckPremission(ctx.Event.UserID, 9) {
			hasArg = true
			err := unbanPlugin(unban)
			if err != nil {
				ctx.SendChain(message.Text(err.Error()))
			} else {
				ctx.SendChain(message.Text(fmt.Sprintf("插件%s全局解禁成功", unban)))
			}
		} else {
			ctx.SendChain(message.Text("你没有权限执行这条指令"))
			return
		}
	}

	// 处理enable指令
	if enable != "" {
		if control.CheckPremission(ctx.Event.UserID, 9) {
			hasArg = true
			err := enablePlugin(enable, ctx.Event.GroupID)
			if err != nil {
				ctx.SendChain(message.Text(err.Error()))
			} else {
				ctx.SendChain(message.Text(fmt.Sprintf("插件%s局部启用成功", enable)))
			}
		} else {
			ctx.SendChain(message.Text("你没有权限执行这条指令"))
			return
		}
	}

	// 处理disable指令
	if disable != "" {
		if control.CheckPremission(ctx.Event.UserID, 9) {
			hasArg = true
			err := disablePlugin(disable, ctx.Event.GroupID)
			if err != nil {
				ctx.SendChain(message.Text(err.Error()))
			} else {
				ctx.SendChain(message.Text(fmt.Sprintf("插件%s局部禁用成功", disable)))
			}
		} else {
			ctx.SendChain(message.Text("你没有权限执行这条指令"))
			return
		}
	}

	// 处理help指令
	if help {
		if control.CheckPremission(ctx.Event.UserID, 5) {
			hasArg = true
			ctx.SendChain(message.Text(usage))
		} else {
			ctx.SendChain(message.Text("你没有权限执行这条指令"))
			return
		}
	}

	// 没有提供参数时，显示帮助
	if !hasArg {
		if control.CheckPremission(ctx.Event.UserID, 5) {
			ctx.SendChain(message.Text(usage))
		} else {
			ctx.SendChain(message.Text("你没有权限执行这条指令"))
			return
		}
	}
}

// getPluginList 获取插件列表
func getPluginList() string {
	// 连接数据库
	db, err := conn.GetDB()
	if err != nil {
		logger.Errorln("数据库连接失败: ", err)
		return "查询插件列表失败: 数据库连接失败"
	}

	var plugins []model.Plugin
	err = db.Find(&plugins).Error
	if err != nil {
		logger.Errorln("查询插件列表失败: ", err)
		return "查询插件列表失败: 数据库查询失败"
	}

	var listBuilder strings.Builder
	listBuilder.WriteString("插件列表:\n")

	for index, plugin := range plugins {
		listBuilder.WriteString(plugin.Name)
		listBuilder.WriteRune('\n')
		listBuilder.WriteRune('\t')
		listBuilder.WriteString("- ")
		listBuilder.WriteString(plugin.Description)
		if index < len(plugins)-1 {
			listBuilder.WriteString("\n")
		}
	}
	return listBuilder.String()
}

// banPlugin 禁用插件
func banPlugin(pluginName string) error {
	// 连接数据库
	db, err := conn.GetDB()
	if err != nil {
		logger.Errorln("连接数据库失败: ", err)
		return errors.New("数据库连接失败")
	}

	// 查询插件是否存在
	var plugin model.GlobalPluginManagement
	err = db.Where("name = ?", pluginName).First(&plugin).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			logger.Errorln("查询插件失败: ", err)
			return errors.New("指令执行失败: 插件不存在")
		}
		logger.Errorln("查询插件失败: ", err)
		return errors.New("指令执行失败: 无法获取插件信息")
	}

	// 判断是否为内建插件
	if plugin.Buitlin {
		return errors.New("指令执行失败: 内建插件不能被禁用")
	}

	// 判断插件是否已经被禁用
	if plugin.IsBaned {
		return errors.New("指令执行失败: 插件已经被禁用")
	}

	// 更新插件状态
	err = db.Model(&plugin).Update("is_baned", true).Error
	if err != nil {
		logger.Errorln("更新插件状态失败: ", err)
		return errors.New("指令执行失败: 无法更新插件状态")
	}

	return nil
}

// unbanPlugin 解禁插件
func unbanPlugin(pluginName string) error {
	// 连接数据库
	db, err := conn.GetDB()
	if err != nil {
		logger.Errorln("连接数据库失败: ", err)
		return errors.New("数据库连接失败")
	}

	// 查询插件是否存在
	var plugin model.GlobalPluginManagement
	err = db.Where("name = ?", pluginName).First(&plugin).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			logger.Errorln("查询插件失败: ", err)
			return errors.New("指令执行失败: 插件不存在")
		}
		logger.Errorln("查询插件失败: ", err)
		return errors.New("指令执行失败: 无法获取插件信息")
	}

	// 判断插件是否为内建插件
	if plugin.Buitlin {
		return errors.New("指令执行失败: 内建插件不能被启用")
	}

	// 判断是否已经被禁用
	if !plugin.IsBaned {
		return errors.New("指令执行失败: 插件已经被全局启用")
	}

	// 更新插件状态
	err = db.Model(&plugin).Update("is_baned", false).Error
	if err != nil {
		logger.Errorln("更新插件状态失败: ", err)
		return errors.New("指令执行失败: 无法更新插件状态")
	}

	return nil
}

// enablePlugin 局部启用插件
func enablePlugin(pluginName string, groupID int64) error {
	// 连接数据库
	db, err := conn.GetDB()
	if err != nil {
		logger.Errorln("连接数据库失败: ", err)
		return errors.New("数据库连接失败")
	}

	// 判断是否为内建插件
	var pluginInfo model.GlobalPluginManagement
	err = db.Where("name = ?", pluginName).First(&pluginInfo).Error
	if err != nil {
		return errors.New("指令执行失败: 查询插件失败")
	}
	if pluginInfo.Buitlin {
		return errors.New("指令执行失败: 内建插件不能被局部启用")
	}

	// 查询插件记录是否存在
	var count int64
	err = db.Model(&model.LocalPluginManagement{}).Count(&count).Error
	if err != nil {
		logger.Errorln("查询插件记录失败: ", err)
		return errors.New("指令执行失败: 查询插件记录失败")
	}

	if count > 0 {
		var plugin model.LocalPluginManagement
		err = db.Where("name = ?", pluginName).First(&plugin).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			logger.Errorln("查询插件失败: ", err)
			return errors.New("指令执行失败: 查询插件失败")
		}

		// 判断群组列表里是否有群组
		if plugin.BanedGroupID != "" {
			// 如果有群组, 判断当前群组是否存在其中
			banedGroups := strings.Split(plugin.BanedGroupID, "|")
			// 如果没被禁用，则索引使用-1表示，否则使用被禁用群的索引
			targetIndex := -1
			for index, group := range banedGroups {
				if strconv.FormatInt(groupID, 10) == group {
					targetIndex = index
					break
				}
			}

			// 当插件被禁用时，如果群已经被禁用，则更新被禁用的群
			if targetIndex != -1 {
				// 构造新的被禁用群列表
				var newBanedGroups = make([]string, len(banedGroups)-1)
				copy(newBanedGroups, banedGroups[:targetIndex])
				copy(newBanedGroups[targetIndex:], banedGroups[targetIndex+1:])
				plugin.BanedGroupID = strings.Join(newBanedGroups, "|")

				// 更新插件状态
				if plugin.BanedGroupID == "" {
					// 如果群组列表为空, 则直接删除记录
					err = db.Where("name = ?", pluginName).Delete(&plugin).Error
					if err != nil {
						logger.Errorln("更新插件状态失败: ", err)
						return errors.New("指令执行失败: 无法更新插件状态")
					}
					return nil
				}

				// 如果群组列表不为空, 则更新记录
				err = db.Model(&plugin).Update("baned_group_id", plugin.BanedGroupID).Error
				if err != nil {
					logger.Errorln("更新插件状态失败: ", err)
					return errors.New("指令执行失败: 无法更新插件状态")
				}
				return nil
			}

			// 在被禁用群列表中没有找到群，则说明群已经被启用，不需要更新
			return errors.New("指令执行失败: 插件已经被局部启用")
		}

		// 当插件记录内容为空时，说明插件已经被全局启用，不需要更新
		// 一般来说不会出现这种情况，因为插件记录内容为空时，该条记录会被删除
		// 但是可能在删除记录时被终止，所以还是写了这个判断
		db.Delete(&plugin)
		return errors.New("指令执行失败: 插件已经被局部启用")
	}

	// 当查件记录不存在时，说明插件已经被全局启用，不需要更新
	return errors.New("指令执行失败: 插件不存在或已被全局启用")
}

// disablePlugin 局部禁用插件
func disablePlugin(pluginName string, groupID int64) error {
	// 连接数据库
	db, err := conn.GetDB()
	if err != nil {
		logger.Errorln("连接数据库失败: ", err)
		return errors.New("数据库连接失败")
	}

	// 判断是否为内建插件
	var pluginInfo model.GlobalPluginManagement
	err = db.Where("name = ?", pluginName).First(&pluginInfo).Error
	if err != nil {
		return errors.New("指令执行失败: 查询插件失败")
	}
	if pluginInfo.Buitlin {
		return errors.New("指令执行失败: 内建插件不能被局部禁用")
	}

	// 查询插件记录是否存在
	var count int64
	err = db.Model(&model.LocalPluginManagement{}).Count(&count).Error
	if err != nil {
		logger.Errorln("查询插件记录失败: ", err)
		return errors.New("指令执行失败: 查询插件记录失败")
	}

	// 判断是否有记录
	if count > 0 {
		// 插件记录存在
		var plugin model.LocalPluginManagement
		err = db.Where("name = ?", pluginName).First(&plugin).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			logger.Errorln("查询插件失败: ", err)
			return errors.New("指令执行失败")
		}

		// 判断群组列表里是否有群组
		if plugin.BanedGroupID != "" {
			// 如果有群组, 判断当前群组是否存在其中
			banedGroups := strings.Split(plugin.BanedGroupID, "|")
			for _, group := range banedGroups {
				if strconv.FormatInt(groupID, 10) == group {
					// 如果群组存在，则不需要更新
					return errors.New("指令执行失败: 插件已经被局部禁用")
				}
			}

			// 如果群组不存在，则添加到群组列表中
			banedGroups = append(banedGroups, strconv.FormatInt(groupID, 10))
			plugin.BanedGroupID = strings.Join(banedGroups, "|")
			// 更新群组列表
			err = db.Model(&plugin).Update("baned_group_id", plugin.BanedGroupID).Error
			if err != nil {
				logger.Errorln("更新插件状态失败: ", err)
				return errors.New("指令执行失败: 无法更新插件状态")
			}
			return nil
		}

		// 如果没有群组, 则直接添加到群组列表中
		plugin.BanedGroupID = strconv.FormatInt(groupID, 10)
		// 更新群组列表
		err = db.Model(&plugin).Update("baned_group_id", plugin.BanedGroupID).Error
		if err != nil {
			logger.Errorln("更新插件状态失败: ", err)
			return errors.New("指令执行失败: 无法更新插件状态")
		}
		return nil
	}

	// 如果没有记录则创建一条记录
	newRecord := model.LocalPluginManagement{
		Name:         pluginName,
		BanedGroupID: strconv.FormatInt(groupID, 10),
	}
	err = db.Create(&newRecord).Error
	if err != nil {
		logger.Errorln("创建插件失败: ", err)
		return errors.New("指令执行失败: 无法创建插件记录")
	}
	return nil
}
