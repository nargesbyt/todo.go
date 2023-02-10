package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewSqlite(dsn string) (*gorm.DB, error) {
	conn, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return conn, nil
}
