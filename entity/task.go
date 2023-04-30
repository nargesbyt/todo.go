package entity

import (
	"database/sql"
	"fmt"
	"github.com/google/jsonapi"
	"time"
)

func (t Task) JSONAPILinks() *jsonapi.Links {
	return &jsonapi.Links{
		"self": fmt.Sprintf("https://localhost:8080/tasks/%d", t.ID),
	}
}
func (t Task) JSONAPIRelationshipLinks(relation string) *jsonapi.Links {
	if relation == "users" {
		return &jsonapi.Links{
			"related": fmt.Sprintf("https://localhost:8080/tasks/%d/users", t.ID),
		}
	}
	return nil
}

type Task struct {
	ID         int64 `gorm:"column:id;primaryKey"`
	Title      string
	Status     string
	CreatedAt  time.Time `gorm:"autoCreateTime"`
	FinishedAt sql.NullTime
	UserID     int64 `gorm:"column:user_id;foreignKey"`
	User       User
}
