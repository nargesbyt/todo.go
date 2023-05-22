package dto

import (
	"github.com/nargesbyt/todo.go/entity"
	"time"
)

type CreateTokenRequest struct {
	Title     string    `json:"title"`
	ExpiredAt time.Time `json:"expired_at"`
}

type UpdateRequest struct {
	Title     string    `json:"title"`
	ExpiredAt time.Time `json:"expired_at"`
	LastUsed  time.Time `json:"last_used"`
	Active    int       `json:"active"`
}

type Tokens struct {
	ID        int64     `jsonapi:"primary,tokens"`
	Title     string    `jsonapi:"attr,title"`
	IssuedAt  time.Time `jsonapi:"attr,issued_at"`
	ExpiredAt time.Time `jsonapi:"attr,expired_at"`
	Token     string    `jsonapi:"attr,token"`
	LastUsed  time.Time `jsonapi:"attr,last_used"`
	Active    int       `jsonapi:"attr,active"`
	User      *User     `jsonapi:"relation,user"`
}

func (r *Tokens) FromEntity(token entity.Token) {
	r.ID = token.ID
	r.Title = token.Title
	r.IssuedAt = token.IssuedAt
	r.Token = token.Token
	r.Active = token.Active

	if token.ExpiredAt.Valid {
		r.ExpiredAt = token.ExpiredAt.Time
	}

	if token.LastUsed.Valid {
		r.LastUsed = token.LastUsed.Time
	}

	user := User{}
	user.FromEntity(token.User)
	r.User = &user
}
