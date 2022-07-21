// Package init Bot初始化相关代码
package init

import (
	logger "github.com/sirupsen/logrus"

	"github.com/DaydreamCafe/Cocoa/V2/src/conn"
)

func init() {
	db, err := conn.GetDB()
	if err != nil {
		logger.Panicln("获取数据库连接失败:", err)
	}
	db.Exec("DROP TABLE IF EXISTS plugin_models;")
}