// Package conn 数据库操作
package conn

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/DaydreamCafe/Cocoa/V2/src/config"
)

var (
	// Config 全局配置
	Config config.Config

	// dsn 数据库连接字符串
	dsn string
)

func init() {
	// load config
	Config.Load()

	// setup dsn
	dsn = "host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai"
	dsn = fmt.Sprintf(
		dsn,
		Config.Database.Address,
		Config.Database.User,
		Config.Database.Password,
		Config.Database.DatabaseName,
		Config.Database.Port,
	)
}

// GetDB 获取数据库对象指针
func GetDB() (*gorm.DB, error) {
	var db *gorm.DB
	// open database
	db, err := gorm.Open(
		postgres.New(
			postgres.Config{
				DSN:                  dsn,
				PreferSimpleProtocol: true, // disables implicit prepared statement usage
			},
		),
		&gorm.Config{},
	)
	return db, err
}
