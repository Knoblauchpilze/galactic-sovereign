package repositories

import (
	"context"
	"fmt"
	"testing"

	"github.com/KnoblauchPilze/user-service/pkg/db"
	"github.com/stretchr/testify/suite"
)

type SqlQueryType int

const (
	QueryBased SqlQueryType = 0
	ExecBased  SqlQueryType = 1
)

type testPoolFunc func(context.Context, db.ConnectionPool) error

type dbPoolInteractionTestCase struct {
	sqlMode SqlQueryType

	handler testPoolFunc

	expectedSql       string
	expectedArguments []interface{}
}

type dbReturnTestFunc func(context.Context, db.ConnectionPool) interface{}

type dbPoolReturnTestCase struct {
	handler         dbReturnTestFunc
	expectedContent interface{}
}

type RepositoryTestSuite struct {
	suite.Suite

	dbPoolInteractionTestCases map[string]dbPoolInteractionTestCase
	dbPoolReturnTestCases      map[string]dbPoolReturnTestCase
}

func (s *RepositoryTestSuite) TestUsesPoolToRunSqlQuery() {
	for name, testCase := range s.dbPoolInteractionTestCases {
		s.T().Run(name, func(t *testing.T) {
			m := &mockConnectionPool{}

			testCase.handler(context.Background(), m)

			called := testCase.getCalledCount(m)
			s.Require().Equal(1, called)
		})
	}
}

func (s *RepositoryTestSuite) TestWhenCallingHandler_ExpectPoolReceivesCorrectSql() {
	for name, testCase := range s.dbPoolInteractionTestCases {
		s.T().Run(name, func(t *testing.T) {
			m := &mockConnectionPool{}

			testCase.handler(context.Background(), m)

			s.Require().Equal(testCase.expectedSql, m.sqlQuery)
		})
	}
}

func (s *RepositoryTestSuite) TestWhenCallingHandler_ExpectPoolReceivesCorrectArguments() {
	for name, testCase := range s.dbPoolInteractionTestCases {
		s.T().Run(name, func(t *testing.T) {
			m := &mockConnectionPool{}

			testCase.handler(context.Background(), m)

			s.Require().Equal(len(testCase.expectedArguments), len(m.args))
			for id, expected := range testCase.expectedArguments {
				actual := m.args[id]
				s.Require().Equal(expected, actual)
			}
		})
	}
}

func (s *RepositoryTestSuite) TestWhenCallingHandler_PropagatesPoolError() {
	for name, testCase := range s.dbPoolInteractionTestCases {
		s.T().Run(name, func(t *testing.T) {
			m := testCase.generateErrorMock(errDefault)

			err := testCase.handler(context.Background(), m)

			s.Require().Equal(errDefault, err)
		})
	}
}

func (s *RepositoryTestSuite) TestWhenRequestSucceeds_ReturnsExpectedValue() {
	for name, testCase := range s.dbPoolReturnTestCases {
		s.T().Run(name, func(t *testing.T) {
			actual := testCase.handler(context.Background(), &mockConnectionPool{})

			s.Require().Equal(testCase.expectedContent, actual)
		})
	}
}

func (tc *dbPoolInteractionTestCase) getCalledCount(m *mockConnectionPool) int {
	switch tc.sqlMode {
	case QueryBased:
		return m.queryCalled
	case ExecBased:
		return m.execCalled
	default:
		panic(fmt.Errorf("Unsupported sql mode %v", tc.sqlMode))
	}
}

func (tc *dbPoolInteractionTestCase) generateErrorMock(err error) *mockConnectionPool {
	switch tc.sqlMode {
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
		panic(fmt.Errorf("Unsupported sql mode %v", tc.sqlMode))
	}
}
