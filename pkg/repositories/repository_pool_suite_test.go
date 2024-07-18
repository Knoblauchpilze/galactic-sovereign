package repositories

import (
	"context"
	"fmt"

	"github.com/stretchr/testify/suite"
)

type RepositoryPoolTestSuite struct {
	suite.Suite

	sqlMode  SqlQueryType
	testFunc testPoolFunc

	expectedSql       string
	expectedArguments []interface{}
}

func (s *RepositoryPoolTestSuite) TestUsesConnectionToRunSqlQuery() {
	mock := &mockConnectionPool{}

	s.testFunc(context.Background(), mock)

	called := s.getCalledCount(mock)
	s.Require().Equal(1, called)
}

func (s *RepositoryPoolTestSuite) TestGeneratesValidSql() {
	mock := &mockConnectionPool{}

	s.testFunc(context.Background(), mock)

	s.Require().Equal(s.expectedSql, mock.sqlQuery)
}

func (s *RepositoryPoolTestSuite) TestProvidesValidArguments() {
	mock := &mockConnectionPool{}

	s.testFunc(context.Background(), mock)

	s.Require().Equal(len(s.expectedArguments), len(mock.args))
	for id, expected := range s.expectedArguments {
		actual := mock.args[id]
		s.Require().Equal(expected, actual)
	}
}

func (s *RepositoryPoolTestSuite) TestPropagatesQueryError() {
	mock := s.generateErrorMock(errDefault)

	err := s.testFunc(context.Background(), mock)

	s.Require().Equal(errDefault, err)
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
