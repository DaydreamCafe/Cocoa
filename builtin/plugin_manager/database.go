// Package pluginmanager 插件管理器
package pluginmanager

import (
	logger "github.com/sirupsen/logrus"

	"github.com/DaydreamCafe/Cocoa/V2/src/conn"
	"github.com/DaydreamCafe/Cocoa/V2/src/model"
)

func initDatabase() {
	// 获取数据库连接
	db, err := conn.GetDB()
	if err != nil {
		logger.Panicln(err)
	}

	// pluginModels 插件元数据表模型切片
	var pluginModels []model.Plugin

	// 更新插件信息
	// 从插件元数据表中获取所有插件信息
	err = db.Find(&pluginModels).Error
	if err != nil {
		logger.Panicln("更新全局插件信息失败：", err)
	}

	/*
		全局插件元数据表初始化
	*/
	// pluginManagementModels 插件元数据表模型切片
	var pluginManagementModels []model.GlobalPluginManagement

	// 从插件管理表中获取所有插件插件管理信息
	err = db.Find(&pluginManagementModels).Error
	if err != nil {
		logger.Panicln("更新全局插件信息失败：", err)
	}

	// exsistedPlugins 已存在的插件ID切片
	var exsistedPlugins []string

	// 删除已经不存在的插件
	for _, pluginManagementModel := range pluginManagementModels {
		flag := true
		for _, pluginModel := range pluginModels {
			if pluginManagementModel.Name == pluginModel.Name {
				flag = false
				exsistedPlugins = append(exsistedPlugins, pluginModel.Name)
				break
			}
		}
		if flag {
			db.Delete(&pluginManagementModel)
		}
	}

	// newPlugins 新插件模型切片
	var newPlugins = make([]model.Plugin, len(pluginModels)-len(exsistedPlugins))

	// 向newPlugins中加入新插件模型
	index := 0
	for _, pluginModel := range pluginModels {
		flag := true
		for _, exsistedPluginName := range exsistedPlugins {
			if pluginModel.Name == exsistedPluginName {
				flag = false
				break
			}
		}
		if flag {
			newPlugins[index] = pluginModel
			index++
		}
	}

	// 向插件管理表中插入新插件模型
	for _, pluginModel := range newPlugins {
		db.Create(&model.GlobalPluginManagement{
			Name:    pluginModel.Name,
			IsBaned: false,
			Buitlin: pluginModel.Buitlin,
		})
	}

	/*
		局部插件元数据表初始化
	*/
	// localPluginManagementModels 局部插件元数据表模型切片
	var localPluginManagementModels []model.LocalPluginManagement
	
	// 从插件管理表中获取所有插件插件管理信息
	err = db.Find(&localPluginManagementModels).Error
	if err != nil {
		logger.Panicln("更新局部插件信息失败：", err)
	}

	// 删除已经不存在的插件
	for _, localPluginManagementModel := range localPluginManagementModels {
		flag := true
		for _, pluginModel := range pluginModels {
			if localPluginManagementModel.Name == pluginModel.Name {
				flag = false
				break
			}
		}
		if flag {
			db.Delete(&localPluginManagementModel)
		}
	}
}
