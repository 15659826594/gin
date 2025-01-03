package connector

import (
	driver "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func GormOpen(config *driver.Config, cfg *gorm.Config) *gorm.DB {
	handle, err := gorm.Open(mysql.Open(config.FormatDSN()), cfg)
	if err != nil {
		return nil
	}
	return handle
}
