package repository

import (
	"awesomeProject/entity"
	"errors"
	"gorm.io/gorm"
	"time"
)

var ErrUserNotFound = errors.New("user not found")

type Users interface {
	Create(email string, password string, username string) (entity.User, error)
	GetUsers(email string, username string) ([]*entity.User, error)
	GetUserByID(userID int64) (entity.User, error)
	GetUserByUsername(username string) (entity.User, error)
	UpdateUsers(id int64, username string, email string, password string) (entity.User, error)
	DeleteUsers(id int64) error
	//UpdatePassword( userID string, password string, tokenHash string) error
}

type users struct {
	db *gorm.DB
}

func NewUser(db *gorm.DB) (Users, error) {
	u := &users{db: db}
	return u, nil
}
func (u *users) Create(email string, password string, username string) (entity.User, error) {
	user := entity.User{
		Email:     email,
		Password:  password,
		Username:  username,
		CreatedAt: time.Now(),
	}
	err := user.HashPassword()
	if err != nil {
		return user, err
	}
	tx := u.db.Create(&user)
	if tx.Error != nil {
		return user, tx.Error
	}
	return user, nil
}
func (u *users) GetUsers(email string, username string) ([]*entity.User, error) {
	var user []*entity.User
	tx := u.db.Where(&entity.User{Email: email, Username: username}).Find(&user)
	if tx.Error != nil {
		if tx.Error == gorm.ErrRecordNotFound {
			return user, ErrUserNotFound
		}
		return user, tx.Error
	}
	return user, nil

}

func (u *users) GetUserByID(userID int64) (entity.User, error) {
	var user entity.User
	tx := u.db.Preload("Tasks").First(&user, userID)
	if tx.Error != nil {
		if tx.Error == gorm.ErrRecordNotFound {
			return user, ErrUserNotFound
		}
		return user, tx.Error
	}
	return user, nil
}
func (u *users) GetUserByUsername(username string) (entity.User, error) {
	var user entity.User
	tx := u.db.First(&user, "username = ?", username)
	if tx.Error != nil {
		if tx.Error == gorm.ErrRecordNotFound {
			return user, ErrUserNotFound
		}
		return user, tx.Error
	}
	return user, nil
}

func (u *users) UpdateUsers(id int64, username string, email string, password string) (entity.User, error) {
	user := entity.User{}
	u.db.First(&user, id)
	user.Password = password
	user.Username = username
	user.Email = email
	tx := u.db.Save(&user)
	if tx.Error != nil {
		return user, tx.Error
	}
	return user, nil
}

func (u *users) DeleteUsers(id int64) error {
	var user entity.User
	tx := u.db.Delete(&user, id)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
