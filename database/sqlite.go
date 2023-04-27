package database

import (
	"awesomeProject/entity"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func NewSqlite(dsn string) (*gorm.DB, error) {
	conn, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	DB = conn
	DB.AutoMigrate(&entity.User{})
	return conn, nil
}

/*func Migrate() {
	DB.AutoMigrate(&entity.User{})
	log.Println("Database migration completed")
}*/
