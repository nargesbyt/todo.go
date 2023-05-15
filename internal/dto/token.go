package dto

import (
	"github.com/nargesbyt/todo.go/entity"
	"time"
)

type CreateTokenRequest struct {
	Title     string `json:"title"`
	ExpiresAt string `json:"expires_at"`
}
type UpdateRequest struct {
	Title     string    `json:"title"`
	ExpiresAt time.Time `json:"expires_at"`
}
type Tokens struct {
	ID          int64     `jsonapi:"primary,tokens"`
	Title       string    `jsonapi:"attr,title"`
	GeneratedAt time.Time `jsonapi:"attr,generated_at"`
	ExpiresAt   time.Time `jsonapi:"attr,expires_at"`
	Token       string    `jsonapi:"attr,token"`
	User        *User     `jsonapi:"relation,user"`
}

func (r *Tokens) FromEntity(token entity.Token) {
	r.ID = token.ID
	r.Title = token.Title
	r.GeneratedAt = token.GeneratedAt
	r.ExpiresAt = token.ExpiresAt
	r.Token = token.Token
	user := User{}
	user.FromEntity(token.User)
	r.User = &user

}
