package task

import (
	"awesomeProject/entity"
	"awesomeProject/handler/user"
	"time"
)

type CreateRequest struct {
	Title string `json:"title"`
}
type UpdateRequest struct {
	Title  string `json:"title"`
	Status string `json:"status"`
}
type Response struct {
	ID         int64          `json:"id" jsonapi:"primary,tasks"`
	Title      string         `json:"title" jsonapi:"attr,title"`
	Status     string         `json:"status" jsonapi:"attr,status"`
	CreatedAt  time.Time      `json:"created_at" jsonapi:"attr,created_at"`
	FinishedAt time.Time      `json:"finished_at" jsonapi:"attr,finished_at"`
	User       *user.Response `jsonapi:"relation,user"`
}

func (r *Response) FromEntity(task entity.Task) {
	r.ID = task.ID
	r.Title = task.Title
	r.Status = task.Status
	r.CreatedAt = task.CreatedAt
	r.FinishedAt = task.FinishedAt.Time
	user := user.Response{}
	user.FromEntity(task.User)
	r.User = &user

}
