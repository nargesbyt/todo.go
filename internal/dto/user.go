package dto

import (
	"github.com/nargesbyt/todo.go/entity"
	"time"
)

type User struct {
	ID        int64     `jsonapi:"primary,users"`
	Username  string    `jsonapi:"attr,username"`
	Email     string    `jsonapi:"attr,email"`
	CreatedAt time.Time `jsonapi:"attr,created_at"`
	UpdatedAt time.Time `jsonapi:"attr,updated_at"`
	Tasks     []*Task   `jsonapi:"relation,tasks"`
}

type UserUpdateRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r *User) FromEntity(user entity.User) {
	r.ID = user.ID
	r.Username = user.Username
	r.Email = user.Email
	r.CreatedAt = user.CreatedAt
	r.UpdatedAt = user.UpdatedAt
	tasks := []*Task{}
	for _, task := range user.Tasks {
		t := Task{}
		t.FromEntity(task)
		tasks = append(tasks, &t)
	}
	r.Tasks = tasks
}
