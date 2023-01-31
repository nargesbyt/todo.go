package entity

import (
	"database/sql"
	"time"
)

type Task struct {
	ID         int64 `gorm:"primaryKey"`
	Title      string
	Status     string
	CreatedAt  time.Time `gorm:"autoCreateTime"`
	FinishedAt sql.NullTime
}
