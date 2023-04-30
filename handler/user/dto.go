package user

import (
	"awesomeProject/entity"
	"awesomeProject/handler/task"
	"time"
)

type UpdateRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
type Response struct {
	ID        int64            `jsonapi:"primary,users"`
	Username  string           `jsonapi:"attr,username"`
	Email     string           ` jsonapi:"attr,email"`
	CreatedAt time.Time        `jsonapi:"attr,created_at"`
	UpdatedAt time.Time        ` jsonapi:"attr,updated_at"`
	Tasks     []*task.Response `jsonapi:"relation,tasks"`
}

func (r *Response) FromEntity(user entity.User) {
	r.ID = user.ID
	r.Username = user.Username
	r.Email = user.Email
	r.CreatedAt = user.CreatedAt
	r.UpdatedAt = user.UpdatedAt
	//tasks := task.Response{}
	//tasks.FromEntity(user.T)
	//r.Tasks
}
