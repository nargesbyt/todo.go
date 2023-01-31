package repository

import (
	"awesomeProject/entity"
	"errors"
	"gorm.io/gorm"
	"time"
)

var ErrTaskNotFound = errors.New("task not found")

type Tasks struct {
	db *gorm.DB
}

/*func (t *Tasks) Init() error {
	sts := `
	CREATE TABLE IF NOT EXISTS tasks
	(id INTEGER PRIMARY KEY, title TEXT NOT NULL , status TEXT NOT NULL DEFAULT "pending",created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, finished_at TIMESTAMP);
	`
	_, err := t.db.Exec(sts)
	if err != nil {
		return err
	}

	return nil
}*/

func NewTasks(db *gorm.DB) (*Tasks, error) {
	t := &Tasks{db: db}

	//err := t.Init()
	/*if err != nil {
		return nil, err
	}*/

	return t, nil
}

func (t *Tasks) Create(title string) (entity.Task, error) {
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

func (t *Tasks) DisplayTask(id int64) (entity.Task, error) {
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

func (t *Tasks) Find(title string, status string) ([]entity.Task, error) {
	var tasks []entity.Task

	tx := t.db.Where(&entity.Task{Title: title, Status: status}).Find(&tasks)

	if tx.Error != nil {
		return tasks, tx.Error
	}
	return tasks, nil
}

func (t *Tasks) Update(id int64, title string, status string) (entity.Task, error) {
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

func (t *Tasks) Delete(id int64) error {
	var task entity.Task
	tx := t.db.Delete(&task, id)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
	//	stm, err := t.db.Prepare("DELETE FROM tasks WHERE id = ?")
	//	if err != nil {
	//		return err
	//	}
	//	defer stm.Close()
	//	res, err := stm.Exec(id)
	//
	//	if err != nil {
	//		return err
	//	}
	//	n, err := res.RowsAffected()
	//	if err != nil {
	//		return err
	//	}
	//	if n == 0 {
	//		return ErrTaskNotFound
	//	}
	//return nil
}
