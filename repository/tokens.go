package repository

import (
	"database/sql"
	"errors"
	"github.com/nargesbyt/todo.go/entity"
	"github.com/nargesbyt/todo.go/internal/random"
	"gorm.io/gorm"
	"time"
)

var ErrTokenNotFound = errors.New("token not found")

type Tokens interface {
	Add(title string, expiredAt time.Time, userId int64) (entity.Token, error)
	Get(id int64) (entity.Token, error)
	List(title string, userId int64) ([]*entity.Token, error)
	GetTokensByUserID(userId int64) ([]*entity.Token, error)
	Update(id int64, title string, expiresAt time.Time, lastUsed time.Time, active int) (entity.Token, error)
	Delete(id int64) error
}

type tokens struct {
	db *gorm.DB
}

func NewTokens(db *gorm.DB) (Tokens, error) {
	token := &tokens{db: db}
	err := db.AutoMigrate(&entity.Token{})
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (t *tokens) Add(title string, expiredAt time.Time, userId int64) (entity.Token, error) {
	randomToken := "todo_pat_"
	randomToken += random.Token(32)

	expireTime := sql.NullTime{}
	err := expireTime.Scan(expiredAt)
	if err != nil {
		return entity.Token{}, err
	}

	token := entity.Token{
		Title:     title,
		ExpiredAt: expireTime,
		Token:     randomToken,
		UserID:    userId,
		Active:    1,
		LastUsed:  expireTime,
	}

	err = token.HashToken()
	if err != nil {
		return entity.Token{}, err
	}

	tx := t.db.Create(&token)
	if tx.Error != nil {
		return entity.Token{}, tx.Error
	}

	token.Token = randomToken

	return token, nil
}

func (t *tokens) Get(id int64) (entity.Token, error) {
	var token entity.Token

	tx := t.db.First(&token, id)
	if tx.Error != nil {
		if tx.Error == gorm.ErrRecordNotFound {
			return token, ErrTokenNotFound
		}

		return token, tx.Error
	}

	return token, nil
}

func (t *tokens) List(title string, userId int64) ([]*entity.Token, error) {
	var tokensList []*entity.Token
	var totalRows int64
	tx := t.db.Model(&entity.Task{Title: title, UserID: userId}).Count(&totalRows)
	if tx.Error != nil {
		return nil, tx.Error
	}
	tx = t.db.Where(&entity.Token{Title: title, UserID: userId}).Find(&tokensList)
	if tx.Error != nil {
		if tx.Error == gorm.ErrRecordNotFound {
			return nil, ErrTokenNotFound
		}
		return nil, tx.Error
	}
	return tokensList, nil

}

func (t *tokens) GetTokensByUserID(userId int64) ([]*entity.Token, error) {
	var tokensList []*entity.Token

	tx := t.db.Where(&entity.Token{UserID: userId}).Find(&tokensList)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return tokensList, nil
}

func (t *tokens) Update(id int64, title string, expiresAt time.Time, lastUsed time.Time, active int) (entity.Token, error) {
	token, err := t.Get(id)
	if err != nil {
		if err == ErrTokenNotFound {
			return token, err
		}
		return token, err
	}

	expireTime := sql.NullTime{}
	err = expireTime.Scan(expiresAt)
	if err != nil {
		return token, err
	}

	lastUsedTime := sql.NullTime{}
	err = expireTime.Scan(lastUsed)
	if err != nil {
		return token, err
	}

	token.ExpiredAt = expireTime
	token.Title = title
	token.LastUsed = lastUsedTime
	token.Active = active

	tx := t.db.Save(&token)
	if tx.Error != nil {
		return token, tx.Error
	}

	return token, nil
}

func (t *tokens) Delete(id int64) error {
	tx := t.db.Delete(&entity.Token{}, id)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}
