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
	CREATE TABLE IF NOT EXISTS permissions
	(
		ID         SERIAL PRIMARY KEY,
		QID        INT NOT NULL,
		LEVEL      INT NOT NULL,
		LAST_LEVEL INT NOT NULL
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
	var queryResult int
	var IfExist bool
	sqlStatement := fmt.Sprintf("SELECT level FROM permissions WHERE qid = %d;", QID)
	err = db.QueryRow(sqlStatement).Scan(&queryResult)
	if err != nil {
		IfExist = false
	} else {
		IfExist = true
	}

	if !IfExist {
		// 如果没有数据则添加权限
		sqlStatement := fmt.Sprintf("INSERT INTO permissions (qid, level, last_level) VALUES (%d, %d, %d);;", QID, Level, defaultPermissionLevel)
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
	var currentLevel int
	var IfExist bool
	sqlStatement := fmt.Sprintf("SELECT level FROM permissions WHERE qid = %d;", QID)
	err = db.QueryRow(sqlStatement).Scan(&currentLevel)
	if err != nil {
		IfExist = false
	} else {
		IfExist = true
	}

	// 修改权限
	if IfExist {
		// 如果已存在则修改
		sqlStatement := fmt.Sprintf("UPDATE permissions SET level = %d, last_level = %d WHERE qid = %d;", Level, currentLevel, QID)
		_, err = db.Exec(sqlStatement)
		if err != nil {
			return err
		}
	} else {
		// 若不存在则添加
		sqlStatement := fmt.Sprintf("INSERT INTO permissions (qid, level, last_level) VALUES (%d, %d, %d);;", QID, Level, defaultPermissionLevel)
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
	var queryResult int
	var IfExist bool
	sqlStatement := fmt.Sprintf("SELECT level FROM permissions WHERE qid = %d;", QID)
	err = db.QueryRow(sqlStatement).Scan(&queryResult)
	if err != nil {
		IfExist = false
	} else {
		IfExist = true
	}

	// 删除用户权限配置
	if IfExist {
		// 如果已存在则修改
		sqlStatement := fmt.Sprintf("DELETE FROM permissions WHERE qid = %d;", QID)
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
	var queryResult int
	var IfExist bool
	sqlStatement := fmt.Sprintf("SELECT level FROM permissions WHERE qid = %d;", QID)
	err = db.QueryRow(sqlStatement).Scan(&queryResult)
	if err != nil {
		IfExist = false
	} else {
		IfExist = true
	}

	// 查询权限
	if IfExist {
		// 如果存在则返回数据库中的权限等级配置值
		var ret int64
		sqlStatement := fmt.Sprintf("SELECT level FROM permissions WHERE qid = %d;", QID)
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

func QueryLastPermission(QID int64) (int64, error) {
	// 连接数据库
	err := database.OpenDB()
	if err != nil {
		return 0, err
	}

	db := database.GetDB()
	defer db.Close()

	// 查询用户权限配置是否已存在
	var queryResult int
	var IfExist bool
	sqlStatement := fmt.Sprintf("SELECT last_level FROM permissions WHERE qid = %d;", QID)
	err = db.QueryRow(sqlStatement).Scan(&queryResult)
	if err != nil {
		IfExist = false
	} else {
		IfExist = true
	}

	// 查询权限
	if IfExist {
		// 如果存在则返回数据库中的权限等级配置值
		var ret int64
		sqlStatement := fmt.Sprintf("SELECT last_level FROM permissions WHERE qid = %d;", QID)
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
	_, err = db.Exec("DROP TABLE IF EXISTS permissions;")
	if err != nil {
		return err
	}

	// 初始化权限表
	sqlStatement := `
	CREATE TABLE IF NOT EXISTS permissions
	(
		ID         SERIAL PRIMARY KEY,
		QID        INT NOT NULL,
		LEVEL      INT NOT NULL,
		LAST_LEVEL INT NOT NULL
	);`
	_, err = db.Exec(sqlStatement)
	if err != nil {
		return err
	}

	return nil
}

func BanUser(QID int64) error {
	isBaned, err := IsBaned(QID)
	if err != nil {
		return err
	}
	if isBaned {
		return errors.New("User already baned!")
	} else {
		return UpdatePermission(QID, -1)
	}

}

func UnbanUser(QID int64) error {
	isBaned, err := IsBaned(QID)
	if err != nil {
		return err
	}
	last_level, err := QueryLastPermission(QID)
	if err != nil {
		return err
	}
	if isBaned {
		return UpdatePermission(QID, last_level)
	} else {
		return errors.New("User not baned!")
	}
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
