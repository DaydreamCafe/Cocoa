// Package control 插件控制模块
package control

import (
	logger "github.com/sirupsen/logrus"
	"gorm.io/gorm"

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
}

var (
	// PluginMetadata 插件元数据
	PluginMetadata *Metadata

	// db 数据库d对象
	db *gorm.DB
)

func init() {
	var err error
	db, err = conn.GetDB()
	if err != nil {
		logger.Panicln("获取数据库连接失败:", err)
	}
	err = db.AutoMigrate(&model.PluginModel{})
	if err != nil {
		logger.Panicln("数据库自动迁移失败:", err)
	}
}

// Registe 向数据库注册插件
func Registe(metadata *Metadata) {
	PluginMetadata = metadata

	result := db.Create(&model.PluginModel{
		Name:    PluginMetadata.Name,
		Version: PluginMetadata.Version,
		Usage:   PluginMetadata.Usage,
	})
	if result.Error != nil {
		logger.Panicln("插件信息写入数据库失败:", result.Error)
	}

	logger.Infof("%s 插件加载成功", metadata.Name)
}
