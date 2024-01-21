package repository

import (
	"github.com/nargesbyt/todo.go/entity"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Get(id int64) (entity.User, error) {
	args := m.Called(id)
	return args.Get(0).(entity.User), args.Error(1)

}
func (m *MockUserRepository) Create(email string, password string, username string) (entity.User, error) {
	args := m.Called(email, password, username)
	return args.Get(0).(entity.User), args.Error(1)

}
func (m *MockUserRepository) Update(id int64, username string, email string, password string) (entity.User, error) {
	args := m.Called(id, username, email, password)
	return args.Get(0).(entity.User), args.Error(1)

}
func (m *MockUserRepository) Delete(id int64) error {
	args := m.Called(id)
	return args.Error(0)

}
func (m *MockUserRepository) Find(email string, username string) ([]*entity.User, error){
	args:= m.Called(email,username)
	return args.Get(0).([]*entity.User), args.Error(1)
}
