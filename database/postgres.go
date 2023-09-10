package database

import (
	"database/sql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgres(db *sql.DB) (*gorm.DB, error) {

	conn, err := gorm.Open(postgres.New(postgres.Config{Conn: db}))
	if err != nil {
		return nil, err
	}
	return conn, nil
}
