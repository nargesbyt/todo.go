package database

import (
	"database/sql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewMysql(db *sql.DB) (*gorm.DB, error) {
	conn, err := gorm.Open(mysql.New(mysql.Config{Conn: db}))
	if err != nil {
		return nil, err
	}
	return conn, nil
}
