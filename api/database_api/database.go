/*
 * 对utils/database的二次封装
 * 便于以后针对插件提供更多方法
 * 原则上插件使用应调用此包
 */
package database_api

import (
	"database/sql"
	"github.com/DaydreamCafe/Cocoa/utils/database"
)

func OpenDB() error {
	return database.OpenDB()
}

func GetDB() *sql.DB {
	return database.GetDB()
}

func CloseDB() error {
	return database.CloseDB()
}
