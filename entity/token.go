package entity

import (
	"database/sql"
	"golang.org/x/crypto/bcrypt"
	"time"
)

const DefaultTokenCost = 14

type Token struct {
	ID        int64 `gorm:"column:id;primaryKey"`
	UserID    int64 `gorm:"column:user_id;foreignKey"`
	Title     string
	Token     string
	IssuedAt  time.Time `gorm:"autoCreateTime"`
	Active    int
	LastUsed  time.Time
	ExpiredAt sql.NullTime
	User      User
}

func (t *Token) HashToken() error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(t.Token), DefaultTokenCost)
	if err != nil {
		return err
	}

	t.Token = string(bytes)

	return nil
}

func (t *Token) VerifyToken(rawToken string) error {
	err := bcrypt.CompareHashAndPassword([]byte(t.Token), []byte(rawToken))
	if err != nil {
		return err
	}

	return nil
}
