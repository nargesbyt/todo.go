package task

import (
	"awesomeProject/entity"
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
	ID         int64     `json:"id"`
	Title      string    `json:"title"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	FinishedAt time.Time `json:"finished_at"`
}

func (r *Response) FromEntity(task entity.Task) {
	r.ID = task.ID
	r.Title = task.Title
	r.Status = task.Status
	r.CreatedAt = task.CreatedAt
	r.FinishedAt = task.FinishedAt.Time

}
