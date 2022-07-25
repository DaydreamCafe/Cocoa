// Package control 插件控制模块
package control

import (
	logger "github.com/sirupsen/logrus"
	zero "github.com/wdvxdr1123/ZeroBot"

	"github.com/DaydreamCafe/Cocoa/V2/src/conn"
	"github.com/DaydreamCafe/Cocoa/V2/src/model"
)

// Metadata 插件元数据
type Metadata struct {
	Name        string
	Version     string
	Description string
	Author      string
	Usage       string
	Buitlin     bool
}

// Registe 向数据库注册插件
func Registe(metadata *Metadata) zero.Engine {
	db, err := conn.GetDB()
	if err != nil {
		logger.Panicln("获取数据库连接失败:", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		logger.Panicln("获取数据库连接失败:", err)
	}
	defer sqlDB.Close()
	
	pluginMetadata := metadata

	result := db.Create(&model.Plugin{
		Name:        pluginMetadata.Name,
		Version:     pluginMetadata.Version,
		Usage:       pluginMetadata.Usage,
		Description: pluginMetadata.Description,
		Buitlin:     pluginMetadata.Buitlin,
	})
	if result.Error != nil {
		logger.Panicln("插件信息写入数据库失败:", result.Error)
	}

	logger.Infof("%s 插件加载成功", metadata.Name)

	return getEngine(*pluginMetadata)
}
