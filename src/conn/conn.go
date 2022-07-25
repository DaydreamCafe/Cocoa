// Package conn 数据库连接
package conn

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"github.com/DaydreamCafe/Cocoa/V2/src/config"
)

var (
	// Config 全局配置
	Config config.Config

	// dsn 数据库连接字符串
	dsn string

	// db 数据库对象指针
	db *gorm.DB
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
	var err error
	// open database
	db, err = gorm.Open(
		postgres.New(
			postgres.Config{
				DSN: dsn,
			},
		),
		&gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
		},
	)
	return db, err
}
