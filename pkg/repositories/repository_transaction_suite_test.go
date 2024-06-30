package repositories

import (
	"context"
	"fmt"

	"github.com/KnoblauchPilze/user-service/pkg/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type testTransactionFunc func(context.Context, db.Transaction) error

type RepositoryTransactionTestSuite struct {
	suite.Suite

	sqlMode  SqlQueryType
	testFunc testTransactionFunc

	expectedSql       []string
	expectedArguments [][]interface{}
}

func (s *RepositoryTransactionTestSuite) TestUsesTransactionToRunSqlQuery() {
	assert := assert.New(s.T())

	mock := &mockTransaction{}

	s.testFunc(context.Background(), mock)

	called := s.getCalledCount(mock)
	assert.Equal(len(s.expectedSql), called)
}

func (s *RepositoryTransactionTestSuite) TestGeneratesValidSql() {
	assert := assert.New(s.T())

	mock := &mockTransaction{}

	s.testFunc(context.Background(), mock)

	assert.Equal(s.expectedSql, mock.sqlQueries)
}

func (s *RepositoryTransactionTestSuite) TestProvidesValidArguments() {
	assert := assert.New(s.T())

	mock := &mockTransaction{}

	s.testFunc(context.Background(), mock)

	assert.Equal(len(s.expectedArguments), len(mock.args))
	for id, expectedArgs := range s.expectedArguments {
		actualArgs := mock.args[id]

		for idArg, expectedArg := range expectedArgs {
			actualArg := actualArgs[idArg]
			assert.Equal(expectedArg, actualArg)
		}
	}
}

func (s *RepositoryTransactionTestSuite) TestPropagatesQueryError() {
	assert := assert.New(s.T())

	mock := s.generateErrorMock(errDefault)

	err := s.testFunc(context.Background(), mock)

	assert.Equal(errDefault, err)
}

func (s *RepositoryTransactionTestSuite) getCalledCount(mock *mockTransaction) int {
	switch s.sqlMode {
	case QueryBased:
		return mock.queryCalled
	case ExecBased:
		return mock.execCalled
	default:
		panic(fmt.Errorf("Unsupported sql mode %v", s.sqlMode))
	}
}

func (s *RepositoryTransactionTestSuite) generateErrorMock(err error) *mockTransaction {
	switch s.sqlMode {
	case QueryBased:
		return &mockTransaction{
			rows: mockRows{
				err: err,
			},
		}
	case ExecBased:
		return &mockTransaction{
			execErr: err,
		}
	default:
		panic(fmt.Errorf("Unsupported sql mode %v", s.sqlMode))
	}
}
