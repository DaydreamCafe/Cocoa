package nickname

import (
	"database/sql"
	"errors"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/DaydreamCafe/Cocoa/api/database_api"
)

func init() {
	// 连接数据库
	db, err := connect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 初始化别名数据表
	sqlStatement := `
	CREATE TABLE IF NOT EXISTS nicknames
	(
    	ID        SERIAL PRIMARY KEY,
    	QID       INT  NOT NULL,
    	USER_NAME TEXT NOT NULL,
    	NICKNAME  TEXT
	);`
	_, err = db.Exec(sqlStatement)
	if err != nil {
		log.Fatal(err)
	}
}

// 连接到数据库的方法
func connect() (*sql.DB, error) {
	err := database_api.OpenDB()
	if err != nil {
		return nil, err
	}
	db := database_api.GetDB()
	return db, nil
}

// 查询记录是否已存在
func queryExistence(db *sql.DB, QID int64) bool {
	var queryResult string
	sqlStatement := fmt.Sprintf("SELECT nickname FROM nicknames WHERE qid = %d;", QID)
	err := db.QueryRow(sqlStatement).Scan(&queryResult)
	if err != nil {
		IfExist := false
		return IfExist
	} else {
		IfExist := true
		return IfExist
	}
}

// AddNickname
// 向数据库种增加一条注释记录的方法
//
// params:
// QID int64 需要增加用户的QQ号
// userName string 该用户QQ名
// nickname string 需要增加的昵称
//
// return:
// error error 错误
func AddNickname(QID int64, userName string, nickname string) error {
	// 连接数据库
	db, err := connect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 查询记录是否已存在
	IfExist := queryExistence(db, QID)
	
	// 添加记录
	if IfExist {
		// 如果存在则返回错误
		return errors.New("Nickname already exists!")
	} else {
		// 如果没有数据则添加权限
		return addNickname(db, QID, userName, nickname)
	}
}

// 操作数据库添加昵称记录的方法
func addNickname(db *sql.DB, QID int64, userName string, nickname string) error {
	// 单独把添加昵称记录的代码封装为一个方法，便于复用
	sqlStatement := fmt.Sprintf("INSERT INTO nicknames (id, qid, user_name, nickname) VALUES (DEFAULT, %d, '%s', '%s');", QID, userName, nickname)
	_, err := db.Exec(sqlStatement)
	if err != nil {
		return err
	}
	return nil
}

// UpdateNickname
// 向数据库种修改一条昵称记录的方法
//
// params:
// QID int64 需要增加用户的QQ号
// nickname string 需要修改的昵称
//
// return:
// error error 错误
func UpdateNickname(QID int64, userName string, nickname string) error {
	// 连接数据库
	db, err := connect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	
	// 添加记录
	// 查询记录是否已存在
	IfExist := queryExistence(db, QID)

	if IfExist {
		// 如果存在则修改记录
		sqlStatement := fmt.Sprintf("UPDATE nicknames SET user_name = '%s', nickname = '%s' WHERE qid = %d;", userName, nickname, QID)
		_, err = db.Exec(sqlStatement)
		if err != nil {
			return err
		}
		return nil
	} else {
		// 如果没有数据则添加
		return addNickname(db, QID, userName, nickname)
	}
}

// RemoveNickname
// 删除昵称记录的方法
//
// params:
// QID int64 需要增加用户的QQ号
//
// return:
// error error 错误
func RemoveNickname(QID int64) error {
	// 连接数据库
	db, err := connect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	
	// 添加记录
	// 查询记录是否已存在
	IfExist := queryExistence(db, QID)

	if IfExist {
		// 如果存在数据则删除
		sqlStatement := fmt.Sprintf("DELETE FROM nicknames WHERE qid = %d;", QID)
		_, err = db.Exec(sqlStatement)
		if err != nil {
			return err
		}
		return nil
	} else {
		// 如果不存在则抛出错误
		return errors.New("Document not exsits!")
	}
}

// QueryNickname
// 查询昵称记录的方法
//
// params:
// QID int64 需要查询用户的QQ号
//
// return:
// nickname string 该用户的昵称
// error error 错误
func QueryNickname(QID int64) (nickname string, err error) {
	// 连接数据库
	db, err := connect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 查询记录
	var queryResult string
	sqlStatement := fmt.Sprintf("SELECT nickname FROM nicknames WHERE qid = %d;", QID)
	err = db.QueryRow(sqlStatement).Scan(&queryResult)
	if err == nil {
		nickname = queryResult
	} else {
		err = errors.New("Queried QID not exsits!")
	}
	return
}

// 重设昵称表
func ResetAllNicknames() error {
	// 连接数据库
	db, err := connect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 删除表
	_, err = db.Exec("DROP TABLE IF EXISTS nicknames;")
	if err != nil {
		return err
	}

	// 初始化权限表
	sqlStatement := `
	CREATE TABLE IF NOT EXISTS nicknames
	(
    	ID        SERIAL PRIMARY KEY,
    	QID       INT  NOT NULL,
    	USER_NAME TEXT NOT NULL,
    	NICKNAME  TEXT
	);`
	_, err = db.Exec(sqlStatement)
	if err != nil {
		return err
	}

	return nil
}
