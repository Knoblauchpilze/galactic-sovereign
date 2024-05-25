package repositories

import (
	"context"
	"fmt"

	"github.com/KnoblauchPilze/user-service/pkg/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type SqlQueryType int

const (
	QueryBased SqlQueryType = 0
	ExecBased  SqlQueryType = 1
)

type testFunc func(context.Context, db.ConnectionPool) error

type RepositoryTestSuite struct {
	suite.Suite

	sqlMode  SqlQueryType
	testFunc testFunc

	expectedSql       string
	expectedArguments []interface{}
}

func (s *RepositoryTestSuite) TestUsesConnectionToRunSqlQuery() {
	assert := assert.New(s.T())

	mock := &mockConnectionPool{}

	s.testFunc(context.Background(), mock)

	called := getCalledCount(s.sqlMode, mock)
	assert.Equal(1, called)
}

func (s *RepositoryTestSuite) TestGeneratesValidSql() {
	assert := assert.New(s.T())

	mock := &mockConnectionPool{}

	s.testFunc(context.Background(), mock)

	assert.Equal(s.expectedSql, mock.sqlQuery)
}

func (s *RepositoryTestSuite) TestProvidesValidArguments() {
	assert := assert.New(s.T())

	mock := &mockConnectionPool{}

	s.testFunc(context.Background(), mock)

	assert.Equal(len(mock.args), len(s.expectedArguments))
	for id, expected := range s.expectedArguments {
		actual := mock.args[id]
		assert.Equal(expected, actual)
	}
}

func (s *RepositoryTestSuite) TestPropagatesQueryError() {
	assert := assert.New(s.T())

	mock := generateErrorMock(s.sqlMode, errDefault)

	err := s.testFunc(context.Background(), mock)

	assert.Equal(errDefault, err)
}

func getCalledCount(mode SqlQueryType, mock *mockConnectionPool) int {
	switch mode {
	case QueryBased:
		return mock.queryCalled
	case ExecBased:
		return mock.execCalled
	default:
		panic(fmt.Errorf("Unsupported sql mode %v", mode))
	}
}

func generateErrorMock(mode SqlQueryType, err error) *mockConnectionPool {
	switch mode {
	case QueryBased:
		return &mockConnectionPool{
			rows: mockRows{
				err: err,
			},
		}
	case ExecBased:
		return &mockConnectionPool{
			execErr: err,
		}
	default:
		panic(fmt.Errorf("Unsupported sql mode %v", mode))
	}
}
