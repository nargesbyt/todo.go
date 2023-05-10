package entity

import "time"

type token struct {
	ID        int64     `json:"ID"`
	UserID    int64     `json:"user-id"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires-at"`
}
