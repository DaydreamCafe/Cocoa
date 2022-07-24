// Package model 数据库模型
package model

// Plugin 插件元数据模型
type Plugin struct {
	ID          uint   `gorm:"primary_key"`
	Name        string `gorm:"column:name"`
	Version     string `gorm:"column:version"`
	Usage       string `gorm:"column:usage"`
	Description string `gorm:"column:description"`
	Buitlin     bool   `gorm:"column:buitlin"`
}

// GlobalPluginManagement 全局插件管理模型
type GlobalPluginManagement struct {
	ID      int    `gorm:"primary_key"`
	Name    string `gorm:"column:name"`
	IsBaned bool   `gorm:"column:is_baned"`
	Buitlin bool   `gorm:"column:buitlin"`
}

// LocalPluginManagement 局部插件管理模型
type LocalPluginManagement struct {
	ID           int    `gorm:"primary_key"`
	Name         string `gorm:"column:name"`
	BanedGroupID string `gorm:"column:baned_group_id"`
}

// UserPremissionModel 用户权限模型
type UserPremissionModel struct {
	ID    int   `gorm:"primary_key"`
	QID   int64 `gorm:"column:qid"`
	Level int64 `gorm:"column:level"`
	IfSU  bool  `gorm:"column:if_su"`
}

// BanedUserModel 禁用用户模型
type BanedUserModel struct {
	ID             int   `gorm:"primary_key"`
	QID            int64 `gorm:"column:qid"`
	BanTimeStamp   int64 `gorm:"column:ban_time_stamp"`
	UnbanTimeStamp int64 `gorm:"column:unban_time_stamp"`
}
