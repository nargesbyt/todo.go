package repository

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/nargesbyt/todo.go/database"
	"github.com/nargesbyt/todo.go/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type AnyTime struct{}

// Match satisfies sqlmock.Argument interface
func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

type TaskSuite struct {
	suite.Suite
	DB    *gorm.DB
	mock  sqlmock.Sqlmock
	tasks Tasks
	users Users
}

func (s *TaskSuite) SetupTest() {
	var (
		db  *sql.DB
		err error
	)
    db, s.mock, err = sqlmock.New()
	s.Require().NoError(err)

	//s.DB,err =database.NewSqlite("todo")
	s.DB, err = database.NewPostgres(db)
	//s.DB, err = gorm.Open("postgres", db)
	require.NoError(s.T(), err)
	//s.DB.Logger.LogMode(true);
	s.tasks, _ = NewTasks(s.DB)
	s.users, _ = NewUsers(s.DB)

}
func (s *TaskSuite) TestGet() {
	expectedTask := entity.Task{
		ID:        1,
		Title:     "New task",
		Status:    "pending",
		CreatedAt: time.Now(),
		UserID: 1,
	}

	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "tasks" WHERE "tasks"."id" = $1 ORDER BY "tasks"."id" LIMIT 1`)).
		WithArgs(expectedTask.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "status", "created_at", "finished_at","user_id"}).
			AddRow(expectedTask.ID, expectedTask.Title, expectedTask.Status, expectedTask.CreatedAt, nil,expectedTask.UserID))

	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE "users"."id" = $1`)).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email"}).
		AddRow(1,"ali","ali@yahoo.com"))
	
	_, err := s.tasks.Get(expectedTask.ID)
	s.Require().NoError(err)
	s.Require().NoError(s.mock.ExpectationsWereMet())
}

func (s *TaskSuite) TestCreate() {
	expectedTask := entity.Task{
		ID: 1,
		Title:  "New task",
		Status: "pending",
		UserID: 1,
	}
	s.mock.ExpectBegin()
	s.mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "tasks" ("title","status","created_at","finished_at","user_id") VALUES ($1,$2,$3,$4,$5) RETURNING "id"`)).
		WithArgs(expectedTask.Title,expectedTask.Status,sqlmock.AnyArg(),sqlmock.AnyArg(),expectedTask.UserID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("1"))

	s.mock.ExpectCommit()
	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE "users"."id" = $1 ORDER BY "users"."id" LIMIT 1`)).
		WithArgs(expectedTask.UserID).WillReturnRows(sqlmock.NewRows([]string{"id","username","email"}).
		AddRow(1,"ali","ali@gmail.com"))

	
	_, err := s.tasks.Create(expectedTask.Title, expectedTask.UserID)
	require.NoError(s.T(), err)
	if err = s.mock.ExpectationsWereMet(); err != nil {
		fmt.Printf("unmet expectation error: %s", err)
	}

}

func (s *TaskSuite) TestFind() {
	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "tasks" WHERE "tasks"."title" = $1 AND "tasks"."status" = $2 AND "tasks"."user_id" = $3 LIMIT 1 OFFSET 2`)).
		WithArgs("New task", "pending", 1).
		WillReturnRows(sqlmock.NewRows([]string{"title", "status", "created_at", "finished_at", "user_Id"}).
			AddRow("New task", "pending", nil, nil, 1))		

	_, err := s.tasks.Find("New task", "pending", 1, 3, 1)
	//assert.Equal(t, expected, task)
	require.NoError(s.T(), err)
	if err = s.mock.ExpectationsWereMet(); err != nil {
		fmt.Printf("unmet expectation error: %s", err)
	}
}

func (s *TaskSuite) TestUpdate() {
	expectedTask := entity.Task{
		ID:        1,
		Title:     "New task",
		Status:    "pending",
		CreatedAt: time.Now(),
	}
	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "tasks" WHERE "tasks"."id" = $1`)).
		WithArgs(expectedTask.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "status", "created_at", "finished_at"}).
			AddRow(expectedTask.ID, expectedTask.Title, expectedTask.Status, expectedTask.CreatedAt, nil))
	s.mock.ExpectBegin()
	s.mock.ExpectExec(regexp.QuoteMeta(`UPDATE "tasks" SET "title"=$1,"status"=$2 WHERE "id" = $3`)).
		WithArgs("updated task", "in progress",expectedTask.ID).WillReturnResult(sqlmock.NewResult(1,1))
	

	s.mock.ExpectCommit()
	task, err := s.tasks.Update(expectedTask.ID, "updated task", "in progress")
	s.Assert().Equal("updated task", task.Title)
	assert.Equal(s.T(),"in progress",task.Status)
	require.NoError(s.T(), err)
	if err = s.mock.ExpectationsWereMet(); err != nil {
		fmt.Printf("unmet expectation error: %s", err)
	}
}

func (s *TaskSuite) TestDelete() {
	s.mock.ExpectBegin()
	s.mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "tasks" WHERE "tasks"."id" = $1`)).
		WithArgs(6).WillReturnResult(sqlmock.NewResult(1,1))
	s.mock.ExpectCommit()
	err := s.tasks.Delete(6)
	require.NoError(s.T(), err)
}

func TestTaskSuite(t *testing.T) {
	suite.Run(t, new(TaskSuite))
}
