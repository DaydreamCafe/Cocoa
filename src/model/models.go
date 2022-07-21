// Package model 数据库模型
package model

// PluginModel 插件元数据模型
type PluginModel struct {
	ID      uint   `gorm:"primary_key"`
	Name    string `gorm:"column:name"`
	Version string `gorm:"column:version"`
	Usage   string `gorm:"column:usage"`
}
