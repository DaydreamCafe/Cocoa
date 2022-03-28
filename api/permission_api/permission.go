/*
 * 权限相关API
 */
package permission_api

import (
	"errors"
	"fmt"
	"io/ioutil"

	log "github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"

	"github.com/DaydreamCafe/Cocoa/utils/database"
)

// 默认权限设置结构体
type PermissionConfig struct {
	MinLevel int64 `yaml:"DefaultLevel"`
}

var permissionConfig PermissionConfig
var defaultPermissionLevel int64

func init() {
	// 读取默认权限设置

	configFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Error(err)
	}
	err = yaml.Unmarshal(configFile, &permissionConfig)
	if err != nil {
		log.Error(err)
	}
	defaultPermissionLevel = permissionConfig.MinLevel

	// 连接数据库
	err = database.OpenDB()
	if err != nil {
		log.Fatal(err)
	}
	db := database.GetDB()
	defer db.Close()

	// 初始化权限表
	log.Debug("开始初始化权限表")
	sqlStatement := `
	CREATE TABLE IF NOT EXISTS permissions(
		ID    SERIAL PRIMARY KEY,
		UID   INT NOT NULL,
		LEVEL INT NOT NULL
	);`
	_, err = db.Exec(sqlStatement)
	if err != nil {
		log.Fatal(err)
	}
}

func AddPermission(QID int64, Level int64) error {
	// 连接数据库
	err := database.OpenDB()
	if err != nil {
		return err
	}

	db := database.GetDB()
	defer db.Close()

	// 查询用户权限配置是否已存在
	var totolResult int
	var existence bool
	sqlStatement := fmt.Sprintf("SELECT level FROM permissions WHERE uid = %d;", QID)
	err = db.QueryRow(sqlStatement).Scan(&totolResult)
	if err != nil {
		existence = false
	} else {
		existence = true
	}

	if !existence {
		// 如果没有数据则添加权限
		sqlStatement := fmt.Sprintf("INSERT INTO permissions (uid, level) VALUES (%d, %d);;", QID, Level)
		_, err = db.Exec(sqlStatement)
		if err != nil {
			return err
		}
	} else {
		// 如果存在则返回错误
		return errors.New("Permission data already exists!")
	}

	return nil
}

func UpdatePermission(QID int64, Level int64) error {
	// 连接数据库
	err := database.OpenDB()
	if err != nil {
		return err
	}

	db := database.GetDB()
	defer db.Close()

	// 查询用户权限配置是否已存在
	var totolResult int
	var existence bool
	sqlStatement := fmt.Sprintf("SELECT level FROM permissions WHERE uid = %d;", QID)
	err = db.QueryRow(sqlStatement).Scan(&totolResult)
	if err != nil {
		existence = false
	} else {
		existence = true
	}

	// 修改权限
	if existence {
		// 如果已存在则修改
		sqlStatement := fmt.Sprintf("UPDATE permissions SET level = %d WHERE uid = %d;", Level, QID)
		_, err = db.Exec(sqlStatement)
		if err != nil {
			return err
		}
	} else {
		// 若不存在则添加
		sqlStatement := fmt.Sprintf("INSERT INTO permissions (uid, level) VALUES (%d, %d);;", QID, Level)
		_, err = db.Exec(sqlStatement)
		if err != nil {
			return err
		}
	}
	return nil
}

func RemovePermission(QID int64) error {
	// 连接数据库
	err := database.OpenDB()
	if err != nil {
		return err
	}

	db := database.GetDB()
	defer db.Close()

	// 查询用户权限配置是否已存在
	var totolResult int
	var existence bool
	sqlStatement := fmt.Sprintf("SELECT level FROM permissions WHERE uid = %d;", QID)
	err = db.QueryRow(sqlStatement).Scan(&totolResult)
	if err != nil {
		existence = false
	} else {
		existence = true
	}

	// 删除用户权限配置
	if existence {
		// 如果已存在则修改
		sqlStatement := fmt.Sprintf("DELETE FROM permissions WHERE uid = %d;", QID)
		_, err = db.Exec(sqlStatement)
		if err != nil {
			return err
		}
	} else {
		// 若不存在则返回错误
		return errors.New("User doesn't exsits!")
	}
	return nil
}

func QueryPermission(QID int64) (int64, error) {
	// 连接数据库
	err := database.OpenDB()
	if err != nil {
		return 0, err
	}

	db := database.GetDB()
	defer db.Close()

	// 查询用户权限配置是否已存在
	var totolResult int
	var existence bool
	sqlStatement := fmt.Sprintf("SELECT level FROM permissions WHERE uid = %d;", QID)
	err = db.QueryRow(sqlStatement).Scan(&totolResult)
	if err != nil {
		existence = false
	} else {
		existence = true
	}

	// 查询权限
	if existence {
		// 如果存在则返回数据库中的权限等级配置值
		var ret int64
		sqlStatement := fmt.Sprintf("SELECT level FROM permissions WHERE uid = %d;", QID)
		err := db.QueryRow(sqlStatement).Scan(&ret)
		if err != nil {
			return 0, err
		}
		return ret, nil
	} else {
		// 若不存在则返回默认值
		return defaultPermissionLevel, nil
	}
}

func CheckPermissions(QID int64, minLevel int64) (bool, error) {
	userLevel, err := QueryPermission(QID)
	if err != nil {
		return false, err
	}
	if userLevel >= minLevel {
		return true, nil
	} else {
		return false, nil
	}
}

func ResetUserLevel(QID int64) error {
	err := UpdatePermission(QID, defaultPermissionLevel)
	return err
}

func RestAllLevel() error {
	// 连接数据库
	err := database.OpenDB()
	if err != nil {
		return err
	}

	db := database.GetDB()
	defer db.Close()

	// 删除表
	_, err = db.Exec("drop table permissions;")
	if err != nil {
		return err
	}

	// 初始化权限表
	sqlStatement := `
	CREATE TABLE IF NOT EXISTS permissions(
		ID    SERIAL PRIMARY KEY,
		UID   INT NOT NULL,
		LEVEL INT NOT NULL
	);`
	_, err = db.Exec(sqlStatement)
	if err != nil {
		return err
	}

	return nil
}

func BanUser(QID int64) error {
	return UpdatePermission(QID, -1)
}

func IsBaned(QID int64) (bool, error) {
	level, err := QueryPermission(QID)
	if err != nil {
		return false, err
	}
	if level == -1 {
		return true, nil
	} else {
		return false, nil
	}
}
