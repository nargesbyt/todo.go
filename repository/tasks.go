package repository

import (
	"errors"
	"github.com/nargesbyt/todo.go/entity"
	"gorm.io/gorm"
	"time"
)

var ErrTaskNotFound = errors.New("task not found")
var ErrUnauthorized = errors.New("permission is denied")

type Tasks interface {
	Create(title string, userId int64) (entity.Task, error)
	Get(id int64) (entity.Task, error)
	Find(title string, status string, userId int64, page int, limit int) ([]*entity.Task, error)
	Update(id int64, title string, status string) (entity.Task, error)
	Delete(id int64) error
}

type tasks struct {
	db *gorm.DB
}

func NewTasks(db *gorm.DB) (Tasks, error) {
	t := &tasks{db: db}

	err := db.AutoMigrate(&entity.Task{})
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (t *tasks) Create(title string, userId int64) (entity.Task, error) {
	user := entity.User{}
	tx := t.db.First(&user, userId)
	if tx.Error != nil {
		return entity.Task{}, tx.Error
	}
	task := entity.Task{
		Title:     title,
		Status:    "pending",
		CreatedAt: time.Now(),
		UserID:    userId,
		User:      user,
	}
	tx = t.db.Create(&task).Preload("User")
	if tx.Error != nil {
		return task, tx.Error
	}
	return task, nil

}

func (t *tasks) Get(id int64) (entity.Task, error) {
	var task entity.Task
	tx := t.db.Preload("User").First(&task, id)
	if tx.Error != nil {
		if tx.Error == gorm.ErrRecordNotFound {
			return task, ErrTaskNotFound
		}
		return task, tx.Error
	}
	return task, nil
}

func (t *tasks) Find(title string, status string, userId int64, page int, limit int) ([]*entity.Task, error) {
	var tasks []*entity.Task
	var totalRows int64
	tx := t.db.Model(&entity.Task{Title: title, Status: status, UserID: userId}).Count(&totalRows)
	if tx.Error != nil {
		return nil, tx.Error
	}
	tx = t.db.Preload("User").Where(&entity.Task{Title: title, Status: status, UserID: userId}).Offset((page - 1) * limit).Limit(limit).Find(&tasks)
	if tx.Error != nil {
		return tasks, tx.Error

	}
	return tasks, nil
}

func (t *tasks) Update(id int64, title string, status string) (entity.Task, error) {
	task := entity.Task{}
	t.db.First(&task, id)
	task.Title = title
	task.Status = status
	tx := t.db.Save(&task)
	if tx.Error != nil {
		return task, nil
	}
	return task, nil

}

func (t *tasks) Delete(id int64) error {
	tx := t.db.Delete(&entity.Task{}, id)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}
