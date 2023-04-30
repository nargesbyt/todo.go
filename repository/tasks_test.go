package repository

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/nargesbyt/todo.go/database"
	"github.com/nargesbyt/todo.go/entity"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"regexp"
	"testing"
	"time"
)

type AnyTime struct{}

// Match satisfies sqlmock.Argument interface
func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

type Suite struct {
	suite.Suite
	DB    *gorm.DB
	mock  sqlmock.Sqlmock
	tasks Tasks
}

func (s *Suite) SetupSuite() {
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
	s.DB, err = database.NewMysql(db)
	//s.DB, err = gorm.Open("postgres", db)
	require.NoError(s.T(), err)
	//s.DB.LogMode(true)
	s.tasks, _ = NewTasks(s.DB)

}
func (s *Suite) TestDisplayTask() {
	expectedTask := entity.Task{
		ID:        11,
		Title:     "New task",
		Status:    "pending",
		CreatedAt: time.Now(),
	}
	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "tasks" WHERE "tasks"." id " = $1 ORDER BY "tasks"." id " LIMIT 1`)).
		WithArgs(expectedTask.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "status", "created_at", "finished_at"}).
			AddRow(expectedTask.ID, expectedTask.Title, expectedTask.Status, expectedTask.CreatedAt, nil))
	_, err := s.tasks.DisplayTask(expectedTask.ID)

	require.NoError(s.T(), err)
	//require.Nil(s.T(), deep.Equal(expectedTask, task))
	if err = s.mock.ExpectationsWereMet(); err != nil {
		fmt.Printf("unmet expectation error: %s", err)
	}

}

func (s *Suite) TestCreate() {
	expectedTask := entity.Task{
		Title:  "New task",
		Status: "pending",
	}
	s.mock.ExpectBegin()
	s.mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "tasks" ("title","status","created_at","finished_at") VALUES ($1,$2,$3,$4) RETURNING " id "`)).
		WithArgs(expectedTask.Title, expectedTask.Status, AnyTime{}, nil).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).
			AddRow(1))
	s.mock.ExpectCommit()
	_, err := s.tasks.Create(expectedTask.Title)
	require.NoError(s.T(), err)
	if err = s.mock.ExpectationsWereMet(); err != nil {
		fmt.Printf("unmet expectation error: %s", err)
	}

}

func (s *Suite) TestFind() {
	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "tasks" WHERE "tasks"."title" = $1 AND "tasks"."status" = $2`)).
		WithArgs("New task", "pending").
		WillReturnRows(sqlmock.NewRows([]string{"title", "status", "created_at", "finished_at"}).
			AddRow("New task", "pending", nil, nil))

	_, err := s.tasks.Find("New task", "pending")
	//assert.Equal(t, expected, task)
	require.NoError(s.T(), err)
	if err = s.mock.ExpectationsWereMet(); err != nil {
		fmt.Printf("unmet expectation error: %s", err)
	}
}

func (s *Suite) TestUpdate() {
	expectedTask := entity.Task{
		ID:        11,
		Title:     "New task",
		Status:    "pending",
		CreatedAt: time.Now(),
	}
	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "tasks" WHERE "tasks"." id " = $1 ORDER BY "tasks"." id " LIMIT 1`)).
		WithArgs(expectedTask.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "status", "created_at", "finished_at"}).
			AddRow(expectedTask.ID, expectedTask.Title, expectedTask.Status, expectedTask.CreatedAt, nil))
	s.mock.ExpectBegin()
	/*s.mock.ExpectQuery(regexp.QuoteMeta(`UPDATE "tasks" SET "tasks"."title"= $1 , "tasks"."status"= $2 WHERE "tasks"."id" = $3`)).
	WithArgs("updated task", "in progress", expectedTask.ID).
	WillReturnRows(sqlmock.NewRows(expectedTask.ID, 1))*/
	s.mock.ExpectCommit()
	_, err := s.tasks.Update(expectedTask.ID, "updated task", "in progress")
	require.NoError(s.T(), err)
	if err = s.mock.ExpectationsWereMet(); err != nil {
		fmt.Printf("unmet expectation error: %s", err)
	}
}

func (s *Suite) TestDelete() {
	s.mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "tasks" WHERE "tasks"."id" = $1`)).
		WithArgs(6)
	err := s.tasks.Delete(6)
	require.NoError(s.T(), err)
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))

}

//func (s *Suite) AfterTest(_, _, string) {
//	require.NoError(s.T(), s.mock.ExpectationsWereMet())
//}
