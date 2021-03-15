package pkg

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"regexp"
	"testing"
)

type user struct {
	ID   int64
	Name string
}

type _DBSuit struct {
	suite.Suite
	mock sqlmock.Sqlmock
	db   *DB
	repo user
}

func (suite *_DBSuit) SetupSuite() {
	conn, mock, _ := sqlmock.New()
	suite.db, _ = NewDBWithMockForTest(false, conn)
	suite.mock = mock
}

func (suite *_DBSuit) TestQueryOne() {
	sql := "SELECT id, name FROM `test` WHERE name = (?)"
	rows := suite.mock.NewRows([]string{"id", "name"}).
		AddRow(1, "TEST")
	suite.mock.
		ExpectQuery(regexp.QuoteMeta(sql)).
		WithArgs(1).
		WillReturnRows(rows)
	var u user
	err := suite.db.QueryOne(&u, sql, 1)
	assert.Nil(suite.T(), err)
	if assert.NotEmpty(suite.T(), u) {
		assert.Equal(suite.T(), int64(1), u.ID)
		assert.Equal(suite.T(), "TEST", u.Name)
	}
}

func (suite *_DBSuit) TestQueryMore() {
	sql := "SELECT id, name FROM `test` WHERE id BETWEEN ? AND ?;"
	id := []int64{1, 2, 3}
	name := []string{"FIRST", "Second", "Third"}
	mockRows := suite.mock.NewRows([]string{"id", "name"}).
		AddRow(id[0], name[0]).
		AddRow(id[1], name[1]).
		AddRow(id[2], name[2])
	suite.mock.ExpectQuery(regexp.QuoteMeta(sql)).
		WithArgs(1, 3).
		WillReturnRows(mockRows)
	rows, err := suite.db.QueryMore(sql, 1, 3)
	assert.Nil(suite.T(), err)
	var index = 0
	for rows.Next() {
		var u user
		err := suite.db.ScanRows(rows, &u)
		assert.Nil(suite.T(), err)
		if assert.NotEmpty(suite.T(), u) {
			suite.Equal(id[index], u.ID)
			suite.Equal(name[index], u.Name)
			index++
		}
	}
}

func TestDbSuite(t *testing.T) {
	suite.Run(t, new(_DBSuit))
}
