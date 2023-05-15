package entity

import (
	"time"
)

type Token struct {
	ID          int64 `gorm:"column:id;primaryKey"`
	Title       string
	Token       string
	GeneratedAt time.Time `gorm:"autoCreateTime"`
	ExpiresAt   time.Time
	UserID      int64 `gorm:"column:user_id;foreignKey"`
	User        User
}
