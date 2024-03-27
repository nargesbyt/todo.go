package repository

import (
	"github.com/nargesbyt/todo.go/entity"
	"github.com/stretchr/testify/mock"
)

type MockTaskRepository struct {
	mock.Mock
}

func (m *MockTaskRepository) Create(title string, userId int64) (entity.Task, error) {
	args := m.Called(title, userId)
	return args.Get(0).(entity.Task), args.Error(1)
}

func (m *MockTaskRepository) Get(id int64) (entity.Task, error) {
	args := m.Called(id)
	return args.Get(0).(entity.Task), args.Error(1)
}

func (m *MockTaskRepository) Find(title string, status string, userId int64, page int, limit int) ([]*entity.Task, error) {
	args := m.Called(title, status, userId, page, limit)
	return args.Get(0).([]*entity.Task), args.Error(1)
}

func (m *MockTaskRepository) Update(id int64, title string, status string) (entity.Task, error) {
	args := m.Called(id, title, status)
	return args.Get(0).(entity.Task), args.Error(1)
}

func (m *MockTaskRepository) Delete(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}
