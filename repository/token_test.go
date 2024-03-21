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
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)
type TokenSuite struct {
	suite.Suite
	DB    *gorm.DB
	mock  sqlmock.Sqlmock
	tokens Tokens
}
func(s *TokenSuite)SetupSuite(){
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
	s.tokens, _ = NewTokens(s.DB)


}
func(s *TokenSuite)TestAdd(){
	s.mock.ExpectBegin()
	s.mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "tokens" ("user_id","title","token","issued_at","active","last_used","expired_at") VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING "id"`)).
		WithArgs(1,"token1",sqlmock.AnyArg(),sqlmock.AnyArg(),1,sqlmock.AnyArg(),sqlmock.AnyArg()).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("1"))
	s.mock.ExpectCommit()

	_,err := s.tokens.Add("token1",time.Now(),1)
	require.NoError(s.T(), err)

	if err = s.mock.ExpectationsWereMet(); err != nil {
		fmt.Printf("unmet expectation error: %s", err)
	}
}
func(s *TokenSuite)TestGet(){
	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "tokens" WHERE "tokens"."id" = $1 ORDER BY "tokens"."id" LIMIT 1`)).
		WithArgs(1).
		WillReturnRows(s.mock.NewRows([]string{"id","user_id","title","token","issued_at","active","last_used","expired_at"}).
		AddRow(1,1,"token1","todo_pat_abc",nil,1,nil,nil))
	_,err := s.tokens.Get(1)
	require.NoError(s.T(), err)
	if err = s.mock.ExpectationsWereMet(); err != nil {
		fmt.Printf("unmet expectation error: %s", err)
	}
}
func(s *TokenSuite)TestGetTokensByUserID(){
	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "tokens" WHERE "tokens"."user_id" = $1`)).
		WithArgs(1).
		WillReturnRows(s.mock.NewRows([]string{"id","user_id","title","token","issued_at","active","last_used","expired_at"}).
		AddRow(1,1,"token1","todo_pat_abc",nil,1,nil,nil))
	_,err := s.tokens.GetTokensByUserID(1)
	require.NoError(s.T(), err)
	if err = s.mock.ExpectationsWereMet(); err != nil {
		fmt.Printf("unmet expectation error: %s", err)
	}
	
}	
func(s *TokenSuite)TestList(){
	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "tokens" WHERE "tokens"."user_id" = $1 AND "tokens"."title" = $2`)).
		WithArgs(1,"token1").
		WillReturnRows(s.mock.NewRows([]string{"id","user_id","title","token","issued_at","active","last_used","expired_at"}).
		AddRow(1,1,"token1","todo_pat_abc",nil,1,nil,nil))
	_,err := s.tokens.List("token1",1)
	require.NoError(s.T(), err)
	if err = s.mock.ExpectationsWereMet(); err != nil {
		fmt.Printf("unmet expectation error: %s", err)
	}
}
func(s *TokenSuite)TestUpdate(){
	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "tokens" WHERE "tokens"."id" = $1`)).
		WithArgs(1).
		WillReturnRows(s.mock.NewRows([]string{"id","user_id","title","token","issued_at","active","last_used","expired_at"}).
		AddRow(1,1,"token1","todo_pat_abc",nil,1,nil,nil))

	s.mock.ExpectBegin()
	s.mock.ExpectExec(regexp.QuoteMeta(`UPDATE "tokens" SET "title"=$1,"active"=$2,"last_used"=$3,"expired_at"=$4 WHERE "id" = $5`)).
		WithArgs("updated token",1,sqlmock.AnyArg(),sqlmock.AnyArg(),1).
		WillReturnResult(sqlmock.NewResult(1,1))
	s.mock.ExpectCommit()

	_, err := s.tokens.Update(1,"updated token",time.Now(),time.Now(),1)
    require.NoError(s.T(), err)
	if err = s.mock.ExpectationsWereMet(); err != nil {
		fmt.Printf("unmet expectation error: %s", err)
	}

}	
func(s *TokenSuite)TestDelete(){
	s.mock.ExpectBegin()
	s.mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "tokens" WHERE "tokens"."id" = $1`)).
		WithArgs(1).WillReturnResult(sqlmock.NewResult(1,1))
	s.mock.ExpectCommit()
	err := s.tokens.Delete(1)
	require.NoError(s.T(), err)
	if err = s.mock.ExpectationsWereMet(); err != nil {
		fmt.Printf("unmet expectation error: %s", err)
	}
}	
func TestTokenSuite(t *testing.T) {
	suite.Run(t, new(TokenSuite))

}