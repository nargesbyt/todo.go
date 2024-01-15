package entity

import (
	"golang.org/x/crypto/bcrypt"
	"time"
)

type User struct {
	ID        int64  `gorm:"column:id;primaryKey"`
	Email     string ` gorm:"unique"`
	Password  string
	Username  string `gorm:"unique"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Tasks     []Task `gorm:"foreignKey:UserID,constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

func (user *User) HashPassword() error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	if err != nil {
		return err
	}
	user.Password = string(bytes)
	return nil
}
func (user *User) CheckPassword(rawPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(rawPassword))
	if err != nil {
		return err
	}
	return nil
}
