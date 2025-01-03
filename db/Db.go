package db

import (
	"database/sql"
	"encoding/json"
	"gin/db/connector"
	"gin/lib/php"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

var mysqlInst = map[string]*sql.DB{}
var gormInst = map[string]*gorm.DB{}

// ConnectSql 数据库初始化，并取得数据库类实例
func ConnectSql(config *mysql.Config, force bool) *sql.DB {
	if config == nil {
		return nil
	}
	if config.Net == "" {
		config.Net = "tcp"
	}
	id := uniqueId(config)
	inst, ok := mysqlInst[id]
	if force || ok {
		return inst
	}
	handle := connector.SqlOpen(config)
	if handle != nil {
		mysqlInst[id] = handle
	}
	return handle
}

// ConnectGorm 数据库初始化，并取得数据库类实例
func ConnectGorm(config *mysql.Config, cfg *gorm.Config, force bool) *gorm.DB {
	if config == nil {
		return nil
	}
	if config.Net == "" {
		config.Net = "tcp"
	}
	id := uniqueId(config)
	inst, ok := gormInst[id]
	if !force && ok {
		return inst
	}
	handle := connector.GormOpen(config, cfg)
	if handle != nil {
		gormInst[id] = handle
	}
	return handle
}

func uniqueId(config *mysql.Config) string {
	str, err := json.Marshal(*config)
	if err != nil {
		return ""
	}
	return php.Md5(string(str))
}
