package repository

import (
	"errors"
	"github.com/nargesbyt/todo.go/entity"
	"gorm.io/gorm"
	"time"
)

var ErrUserNotFound = errors.New("user not found")

type Users interface {
	Create(email string, password string, username string) (entity.User, error)
	Find(email string, username string) ([]*entity.User, error)
	Get(userID int64) (entity.User, error)
	GetUserByUsername(username string) (entity.User, error)
	Update(id int64, username string, email string, password string) (entity.User, error)
	Delete(id int64) error
	//UpdatePassword( userID string, password string, tokenHash string) error
}

type users struct {
	db *gorm.DB
}

func NewUsers(db *gorm.DB) (Users, error) {
	u := &users{db: db}

	err := db.AutoMigrate(&entity.User{})
	if err != nil {
		return nil, err
	}

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
func (u *users) Find(email string, username string) ([]*entity.User, error) {
	var user []*entity.User
	tx := u.db.Preload("Tasks").Where(&entity.User{Email: email, Username: username}).Find(&user)
	if tx.Error != nil {
		if tx.Error == gorm.ErrRecordNotFound {
			return user, ErrUserNotFound
		}
		return user, tx.Error
	}
	return user, nil

}

func (u *users) Get(userID int64) (entity.User, error) {
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
	tx := u.db.Preload("Tasks").First(&user, "username = ?", username)
	if tx.Error != nil {
		if tx.Error == gorm.ErrRecordNotFound {
			return user, ErrUserNotFound
		}
		return user, tx.Error
	}
	return user, nil
}

func (u *users) Update(id int64, username string, email string, password string) (entity.User, error) {
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

func (u *users) Delete(id int64) error {
	tx := u.db.Delete(&entity.User{}, id)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}
