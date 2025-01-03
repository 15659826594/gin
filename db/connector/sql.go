package connector

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	driver "github.com/go-sql-driver/mysql"
)

func SqlOpen(config *driver.Config) *sql.DB {
	handle, err := sql.Open("mysql", config.FormatDSN())
	if err != nil {
		return nil
	}
	err = handle.Ping()
	if err != nil {
		return nil
	}
	return handle
}
