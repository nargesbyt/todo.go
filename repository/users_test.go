package repository

import (
	"database/sql"
	"time"
	//"database/sql/driver"
	"fmt"
	"regexp"
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/nargesbyt/todo.go/database"
	"github.com/nargesbyt/todo.go/entity"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)
 
type UserSuite struct {
	suite.Suite
	DB    *gorm.DB
	mock  sqlmock.Sqlmock
	users Users
}
func(s *UserSuite)SetupSuite(){
	var (
		db  *sql.DB
		err error
	)
    db, s.mock, err = sqlmock.New()
	if err != nil {
		fmt.Println("error creating mock database")
		return
	}
	require.NoError(s.T(), err)
	s.DB, err = database.NewPostgres(db)
	require.NoError(s.T(), err)
	s.users, _ = NewUsers(s.DB)

}
func (s *UserSuite) TestCreate() {
	expectedUser:=entity.User{
		ID: 1,
		Username: "ali",
		Email: "ali@gmail.com",

	}
	s.mock.ExpectBegin()
	s.mock.ExpectQuery("INSERT INTO \"users\" (.+) VALUES (.+) RETURNING \"id\"").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("1"))
	s.mock.ExpectCommit()
	_,err := s.users.Create(expectedUser.Email,expectedUser.Password,expectedUser.Username)

	require.NoError(s.T(), err)
	//s.Assert().Equal(s.T(),user.Password,expectedUser.Password)
	if err = s.mock.ExpectationsWereMet(); err != nil {
		fmt.Printf("unmet expectation error: %s", err)
	}

}
func(s *UserSuite) TestGetUserByID(){
	expectedUser:= entity.User{
		ID: 1,
		Username: "ali",
		Email: "ali@gmail.com",
		Password: "123",
		CreatedAt: time.Now(),

	}
	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE "users"."id" = $1 ORDER BY "users"."id" LIMIT 1`)).
		WithArgs(expectedUser.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "created_at", "updated_at"}).
			AddRow(expectedUser.ID,expectedUser.Username,expectedUser.Email,expectedUser.CreatedAt,nil))

	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "tasks" WHERE "tasks"."user_id" = $1`)).
		WithArgs(expectedUser.ID).WillReturnRows(sqlmock.NewRows([]string{"id", "title", "status", "created_at", "finished_at","user_id"}).
			AddRow(1,"task1","pending",time.Now(), nil,1))

	_,err := s.users.GetUserByID(expectedUser.ID)
	require.NoError(s.T(), err)
	if err = s.mock.ExpectationsWereMet(); err != nil {
		fmt.Printf("unmet expectation error: %s", err)
	}

}
func(s *UserSuite) TestGetUsers(){
	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE "users"."email" = $1 AND "users"."username" = $2`)).
		WithArgs("ali@gmail.com","ali").
		WillReturnRows(s.mock.NewRows([]string{"id", "username", "email", "created_at", "updated_at"}).
			AddRow(1,"ali","ali@gmail.com",time.Now(),nil))

	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "tasks" WHERE "tasks"."user_id" = $1`)).
		WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id", "title", "status", "created_at", "finished_at","user_id"}).
			AddRow(1,"task1","pending",time.Now(), nil,1))

	_,err := s.users.GetUsers("ali@gmail.com","ali")
	require.NoError(s.T(), err)
	if err = s.mock.ExpectationsWereMet(); err != nil {
		fmt.Printf("unmet expectation error: %s", err)
	}
	
}	
func (s *UserSuite) TestGetUserByUsername(){
	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE username = $1 ORDER BY "users"."id" LIMIT 1`)).
		WithArgs("ali").
		WillReturnRows(s.mock.NewRows([]string{"id", "username", "email", "created_at", "updated_at"}).
			AddRow(1,"ali","ali@gmail.com",time.Now(),nil))

	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "tasks" WHERE "tasks"."user_id" = $1`)).
		WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id", "title", "status", "created_at", "finished_at","user_id"}).
			AddRow(1,"task1","pending",time.Now(), nil,1))

	_,err := s.users.GetUserByUsername("ali")
	require.NoError(s.T(), err)
	if err = s.mock.ExpectationsWereMet(); err != nil {
		fmt.Printf("unmet expectation error: %s", err)
	}

}
func (s *UserSuite) TestUpdateUsers(){
	expectedUser:= entity.User{
		ID: 1,
		Username: "ali",
		Email: "ali@gmail.com",
		Password: "123",
		CreatedAt: time.Now(),

	}
	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE "users"."id" = $1`)).
		WithArgs(expectedUser.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "password", "email", "created_at","updated_at"}).
			AddRow(expectedUser.ID,expectedUser.Username,expectedUser.Password,expectedUser.Email,expectedUser.CreatedAt,nil))

	s.mock.ExpectBegin()
	s.mock.ExpectExec(`UPDATE "users" SET .+ WHERE .+`).
		WithArgs("john@yahoo.com","abc","john",sqlmock.AnyArg(),1).WillReturnResult(sqlmock.NewResult(1,1))

	s.mock.ExpectCommit()

	_,err := s.users.UpdateUsers(1,"john","john@yahoo.com","abc")	
	require.NoError(s.T(), err)
	if err = s.mock.ExpectationsWereMet(); err != nil {
		fmt.Printf("unmet expectation error: %s", err)
	}
}


func (s *UserSuite) TestDeleteUsers(){
	s.mock.ExpectBegin()
	s.mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "users" WHERE "users"."id" = $1`)).
		WithArgs(1).WillReturnResult(sqlmock.NewResult(1,1))
	s.mock.ExpectCommit()
	err := s.users.DeleteUsers(1)
	require.NoError(s.T(), err)
}
func TestUserSuite(t *testing.T) {
	suite.Run(t, new(UserSuite))

}