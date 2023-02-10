package repository

import (
	"awesomeProject/entity"
	"errors"
	"gorm.io/gorm"
	"time"
)

var ErrTaskNotFound = errors.New("task not found")

type Tasks interface {
	Create(title string) (entity.Task, error)
	DisplayTask(id int64) (entity.Task, error)
	Find(title string, status string) ([]entity.Task, error)
	Update(id int64, title string, status string) (entity.Task, error)
	Delete(id int64) error
}

type tasks struct {
	db *gorm.DB
}

func NewTasks(db *gorm.DB) (Tasks, error) {
	t := &tasks{db: db}

	//err := t.Init()
	/*if err != nil {
		return nil, err
	}*/

	return t, nil
}

func (t *tasks) Create(title string) (entity.Task, error) {
	task := entity.Task{
		Title:     title,
		Status:    "pending",
		CreatedAt: time.Now(),
	}
	tx := t.db.Create(&task)
	if tx.Error != nil {
		return task, tx.Error
	}
	return task, nil

}

func (t *tasks) DisplayTask(id int64) (entity.Task, error) {
	var task entity.Task
	tx := t.db.First(&task, id)
	if tx.Error != nil {
		if tx.Error == gorm.ErrRecordNotFound {
			return task, ErrTaskNotFound
		}
		return task, tx.Error
	}
	return task, nil

}

func (t *tasks) Find(title string, status string) ([]entity.Task, error) {
	var tasks []entity.Task

	tx := t.db.Where(&entity.Task{Title: title, Status: status}).Find(&tasks)

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
	var task entity.Task
	tx := t.db.Delete(&task, id)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
