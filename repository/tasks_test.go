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

type Suite struct {
	suite.Suite
	DB    *gorm.DB
	mock  sqlmock.Sqlmock
	tasks Tasks
	//users Users
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
	//s.mock.ExpectQuery("SELECT VERSION()").WillReturnRows(sqlmock.NewRows([]string{"version"}).AddRow("5.7.0"))
	if db == nil {
		fmt.Println("mock db is null")
	}
	if s.mock == nil {
		fmt.Println("sqlmock is null")
	}
	require.NoError(s.T(), err)
	s.DB, err = database.NewPostgres(db)
	//s.DB, err = gorm.Open("postgres", db)
	require.NoError(s.T(), err)
	//s.DB.LogMode(true)
	s.tasks, _ = NewTasks(s.DB)
	//s.users, _ = NewUsers(s.DB)
}
func (s *Suite) TestGet() {
	expectedTask := entity.Task{
		ID:        1,
		Title:     "New task",
		Status:    "pending",
		CreatedAt: time.Now(),
		UserID:    8,
	}
	expectedUser := entity.User{
		ID:        8,
		Username:  "ali",
		Email:     "ali@gmail.com",
		CreatedAt: time.Now(),
	} 
	
	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "tasks" WHERE "tasks"." id " = $1 ORDER BY "tasks"." id " LIMIT 1 `)).
		WithArgs(expectedTask.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "status", "created_at", "finished_at"}).
			AddRow(expectedTask.ID, expectedTask.Title, expectedTask.Status, expectedTask.CreatedAt, nil))

	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE "users"."user_id" IN ($1)`)).
		WithArgs(expectedTask.UserID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "created_at", "updated_at"}).
			AddRow(expectedUser.ID, expectedUser.Username, expectedUser.Email, expectedUser.CreatedAt, nil))
	
	//s.mock.MatchExpectationsInOrder(true)
	/*s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT "tasks"."id","tasks"."title","tasks"."status","tasks"."created_at","tasks"."finished_at","tasks"."user_id" FROM "tasks" JOIN "users" ON users.id = tasks.user_id WHERE "tasks"."id" = $1 ORDER BY "tasks"."id" LIMIT 1"` )).
		WithArgs(expectedTask.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "status", "created_at", "finished_at","user_id","username","email","created_at","updated_at"}).
			AddRow(expectedTask.ID, expectedTask.Title, expectedTask.Status, expectedTask.CreatedAt, nil,expectedTask.UserID,expectedUser.Username, expectedUser.Email, expectedUser.CreatedAt, nil))*/

	_, err := s.tasks.Get(expectedTask.ID)

	require.NoError(s.T(), err)
	//require.Nil(s.T(), deep.Equal(expectedTask, task))
	if err = s.mock.ExpectationsWereMet(); err != nil {
		fmt.Printf("there were unfulfilled expectation : %s", err)
	}

}

/*func (s *Suite) TestCreate() {
	var UserID=2

	//user:= entity.User{}
	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE "users"."id" = $1 ORDER BY "users"."id" LIMIT 1`)).
		WithArgs(UserID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "email", "username"}).
			AddRow(2, "www.alireza@gmail.com","ali"))
	/*user1:= entity.User{
		ID: 2,
		Username: "ali",
		Email: "www.alireza@gmail.com",
	}*/

	/*expectedTask := entity.Task{
		Title:  "New task",
		Status: "pending",
		CreatedAt: time.Now(),
		UserID: 2,
		//User: user1,
	}

	s.mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "tasks" ("title","status","created_at","finished_at","user_id") VALUES ($1,$2,$3,$4,$5) RETURNING " id "`)).
		WithArgs(expectedTask.Title, expectedTask.Status, AnyTime{}, nil, expectedTask.UserID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "status", "created_at", "finished_at","user_id","username", "email", "created_at", "updated_at"}).
		AddRow(1,expectedTask.Title,expectedTask.Status,expectedTask.CreatedAt,nil,expectedTask.UserID,"ali","alireza@gmail.com",nil,nil))

	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE "users"."user_id" IN ($1)`)).
		WithArgs(expectedTask.UserID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "created_at", "updated_at"}).
			AddRow( expectedTask.UserID,"ali","ali@gmail.com",time.Now(),nil))

	//s.mock.ExpectBegin()

	//s.mock.ExpectCommit()


	s.mock.MatchExpectationsInOrder(true)
	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE "users"."user_id" IN ($1)`)).
		WithArgs(expectedTask.UserID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "created_at", "updated_at"}).
			AddRow( expectedTask.UserID,"ali","ali@gmail.com",time.Now(),nil))


	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT "tasks"."id","tasks"."title","tasks"."status","tasks"."created_at","tasks"."finished_at","tasks"."user_id","users"."username","users"."email","users"."created_at","users"."updated_at", FROM "tasks" LEFT JOIN "users" ON "tasks"."user_id"="users"."user_id"  WHERE "userss"."user_id" IN ($1)`)).
		WithArgs(expectedTask.UserID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "created_at", "updated_at"}).
			AddRow( expectedTask.UserID,"ali","ali@gmail.com",time.Now(),nil))

	_, err := s.tasks.Create(expectedTask.Title, expectedTask.UserID)
	require.NoError(s.T(), err)
	if err = s.mock.ExpectationsWereMet(); err != nil {
		fmt.Printf("unmet expectation error: %s", err)
	}

}*/

/*func (s *Suite) TestFind() {
	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "tasks" WHERE "tasks"."title" = $1 AND "tasks"."status" = $2 AND "tasks"."user_id" = $3`)).
	WithArgs("task1","pending",5).
	WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))
	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE "users"."user_id" IN ($1)`)).
		WithArgs(5).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "created_at", "updated_at"}).
			AddRow( 5,"ali","ali@gmail.com",time.Now(),nil))


	_, err := s.tasks.Find("task1", "pending", 5, 1, 3)
	//assert.Equal(t, expected, task)
	require.NoError(s.T(), err)
	if err = s.mock.ExpectationsWereMet(); err != nil {
		fmt.Printf("unmet expectation error: %s", err)
	}
}*/

func (s *Suite) TestUpdate() {
	expectedTask := entity.Task{
		ID:        11,
		Title:     "New task",
		Status:    "pending",
		CreatedAt: time.Now(),
		UserID:    2,
	}
	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "tasks" WHERE "tasks"." id " = $1 ORDER BY "tasks"." id " LIMIT 1`)).
		WithArgs(expectedTask.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "status", "created_at", "finished_at", "user_id"}).
			AddRow(expectedTask.ID, expectedTask.Title, expectedTask.Status, expectedTask.CreatedAt, nil, expectedTask.UserID))
	s.mock.ExpectBegin()
	s.mock.ExpectQuery(regexp.QuoteMeta(`UPDATE "tasks" SET "tasks"."title"= $1 , "tasks"."status"= $2 WHERE "tasks"."id" = $3`)).
		WithArgs("updated task", "in progress", expectedTask.ID).
		WillReturnRows(sqlmock.NewRows([]string{"title", "status"}).
			AddRow("updated task", "in progress"))
	s.mock.ExpectCommit()
	_, err := s.tasks.Update(expectedTask.ID, "updated task", "in progress")
	require.NoError(s.T(), err)
	if err = s.mock.ExpectationsWereMet(); err != nil {
		fmt.Printf("unmet expectation error: %s", err)
	}
}

func (s *Suite) TestDelete() {
	s.mock.ExpectBegin()
	s.mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "tasks" WHERE "tasks"."id" = $1`)).
		WithArgs(6).WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()
	err := s.tasks.Delete(6)
	require.NoError(s.T(), err)
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))

}

//func (s *Suite) AfterTest(_, _, string) {
//	require.NoError(s.T(), s.mock.ExpectationsWereMet())
//}
