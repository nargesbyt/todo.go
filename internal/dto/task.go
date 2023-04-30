package dto

import (
	"github.com/nargesbyt/todo.go/entity"
	"time"
)

type TaskCreateRequest struct {
	Title string `json:"title"`
}
type TaskUpdateRequest struct {
	Title  string `json:"title"`
	Status string `json:"status"`
}
type Task struct {
	ID         int64     `jsonapi:"primary,tasks"`
	Title      string    `jsonapi:"attr,title"`
	Status     string    `jsonapi:"attr,status"`
	CreatedAt  time.Time `jsonapi:"attr,created_at"`
	FinishedAt time.Time `jsonapi:"attr,finished_at"`
	User       *User     `jsonapi:"relation,user"`
}

func (r *Task) FromEntity(task entity.Task) {
	r.ID = task.ID
	r.Title = task.Title
	r.Status = task.Status
	r.CreatedAt = task.CreatedAt
	r.FinishedAt = task.FinishedAt.Time
	user := User{}
	user.FromEntity(task.User)
	r.User = &user
}
