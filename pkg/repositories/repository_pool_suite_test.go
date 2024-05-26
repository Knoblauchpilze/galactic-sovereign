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

type testPoolFunc func(context.Context, db.ConnectionPool) error

type RepositoryPoolTestSuite struct {
	suite.Suite

	sqlMode  SqlQueryType
	testFunc testPoolFunc

	expectedSql       string
	expectedArguments []interface{}
}

func (s *RepositoryPoolTestSuite) TestUsesConnectionToRunSqlQuery() {
	assert := assert.New(s.T())

	mock := &mockConnectionPool{}

	s.testFunc(context.Background(), mock)

	called := s.getCalledCount(mock)
	assert.Equal(1, called)
}

func (s *RepositoryPoolTestSuite) TestGeneratesValidSql() {
	assert := assert.New(s.T())

	mock := &mockConnectionPool{}

	s.testFunc(context.Background(), mock)

	assert.Equal(s.expectedSql, mock.sqlQuery)
}

func (s *RepositoryPoolTestSuite) TestProvidesValidArguments() {
	assert := assert.New(s.T())

	mock := &mockConnectionPool{}

	s.testFunc(context.Background(), mock)

	assert.Equal(len(s.expectedArguments), len(mock.args))
	for id, expected := range s.expectedArguments {
		actual := mock.args[id]
		assert.Equal(expected, actual)
	}
}

func (s *RepositoryPoolTestSuite) TestPropagatesQueryError() {
	assert := assert.New(s.T())

	mock := s.generateErrorMock(errDefault)

	err := s.testFunc(context.Background(), mock)

	assert.Equal(errDefault, err)
}

func (s *RepositoryPoolTestSuite) getCalledCount(mock *mockConnectionPool) int {
	switch s.sqlMode {
	case QueryBased:
		return mock.queryCalled
	case ExecBased:
		return mock.execCalled
	default:
		panic(fmt.Errorf("Unsupported sql mode %v", s.sqlMode))
	}
}

func (s *RepositoryPoolTestSuite) generateErrorMock(err error) *mockConnectionPool {
	switch s.sqlMode {
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
		panic(fmt.Errorf("Unsupported sql mode %v", s.sqlMode))
	}
}
