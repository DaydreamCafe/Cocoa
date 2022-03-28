/*
 * 基础数据库模块
 * 如果是插件使用，请调用api/database_api而不是直接导入本包
 */
package database

import (
	"database/sql"
	"fmt"
	"io/ioutil"

	log "github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"

	_ "github.com/lib/pq"
)

// 定义数据库配置结构体
type Config struct {
	DatabaseConfig struct {
		Address      string `yaml:"Address"`
		Port         int    `yaml:"Port"`
		User         string `yaml:"User"`
		Password     string `yaml:"Password"`
		DatabaseName string `yaml:"DatabaseName"`
	} `yaml:"Database"`
}

/*
 * DB SSL Mode
 * disable     - No SSL
 * require     - Always SSL (skip verification)
 * verify-ca   - Always SSL (verify that the certificate presented by the server was signed by a trusted CA)
 * verify-full - Always SSL (verify that the certification presented by the server was signed by a trusted
   CA and the server host name matches the one in the certificate)
 * See more in https://pkg.go.dev/github.com/lib/pq@v1.10.4#section-readme
*/
const dbSSLMode = "disable"

var config Config
var db *sql.DB
var connStr string

func init() {
	// 读取配置文件
	dbConfigFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Error(err)
	}
	err = yaml.Unmarshal(dbConfigFile, &config)
	if err != nil {
		log.Error(err)
	}
	databaseConfig := config.DatabaseConfig

	// 数据库连接配置
	connStr = fmt.Sprintf(
		"user=%s password=%s host=%s port=%d dbname=%s sslmode=%s",
		databaseConfig.User,
		databaseConfig.Password,
		databaseConfig.Address,
		databaseConfig.Port,
		databaseConfig.DatabaseName,
		dbSSLMode,
	)

}

func OpenDB() error {
	var err error
	// 连接到数据库
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		return err
	}

	// 验证连接
	if err = db.Ping(); err != nil {
		return err
	}

	return err
}

func GetDB() *sql.DB {
	return db
}

func CloseDB() error {
	return db.Close()
}
