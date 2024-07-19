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

type dbPoolInteractionTestCase struct {
	sqlMode SqlQueryType

	handler testPoolFunc

	expectedSql       string
	expectedArguments []interface{}
}

type testPoolFunc func(context.Context, db.ConnectionPool) error

type dbPoolReturnTestCase struct {
	handler         testPoolReturnFunc
	expectedContent interface{}
}

type testPoolReturnFunc func(context.Context, db.ConnectionPool) interface{}

type dbPoolSingleValueTestCase struct {
	handler testPoolFunc

	expectedGetSingleValueCalls int
	expectedScanCalls           int
	expectedScannedProps        [][]interface{}
}

type dbPoolGetAllTestCase struct {
	handler testPoolFunc

	expectedGetAllCalls  int
	expectedScanCalls    int
	expectedScannedProps [][]interface{}
}

type RepositoryTestSuite struct {
	suite.Suite

	dbPoolInteractionTestCases map[string]dbPoolInteractionTestCase
	dbPoolSingleValueTestCases map[string]dbPoolSingleValueTestCase
	dbPoolGetAllTestCases      map[string]dbPoolGetAllTestCase
	dbPoolReturnTestCases      map[string]dbPoolReturnTestCase
}

func (s *RepositoryTestSuite) TestPool_ExpectCorrectNumberOfCalls() {
	for name, testCase := range s.dbPoolInteractionTestCases {
		s.T().Run(name, func(t *testing.T) {
			m := &mockConnectionPoolNew{}

			testCase.handler(context.Background(), m)

			called := testCase.getCalledCount(m)
			s.Require().Equal(1, called)
		})
	}
}

func (s *RepositoryTestSuite) TestPool_ExpectCorrectSqlQuery() {
	for name, testCase := range s.dbPoolInteractionTestCases {
		s.T().Run(name, func(t *testing.T) {
			m := &mockConnectionPoolNew{}

			testCase.handler(context.Background(), m)

			s.Require().Equal(testCase.expectedSql, m.sqlQuery)
		})
	}
}

func (s *RepositoryTestSuite) TestPool_ExpectCorrectSqlArguments() {
	for name, testCase := range s.dbPoolInteractionTestCases {
		s.T().Run(name, func(t *testing.T) {
			m := &mockConnectionPoolNew{}

			testCase.handler(context.Background(), m)

			s.Require().Equal(len(testCase.expectedArguments), len(m.args))
			for id, expected := range testCase.expectedArguments {
				actual := m.args[id]
				s.Require().Equal(expected, actual)
			}
		})
	}
}

func (s *RepositoryTestSuite) TestPool_PropagatesError() {
	for name, testCase := range s.dbPoolInteractionTestCases {
		s.T().Run(name, func(t *testing.T) {
			m := testCase.generateErrorMock(errDefault)

			err := testCase.handler(context.Background(), m)

			s.Require().Equal(errDefault, err)
		})
	}
}

func (s *RepositoryTestSuite) TestGetSingleValue_ExpectCorrectNumberOfCalls() {
	for name, testCase := range s.dbPoolSingleValueTestCases {
		s.T().Run(name, func(t *testing.T) {
			m := &mockConnectionPoolNew{}

			testCase.handler(context.Background(), m)

			s.Require().Equal(testCase.expectedGetSingleValueCalls, m.rows.getSingleValueCalled)
		})
	}
}

func (s *RepositoryTestSuite) TestGetSingleValue_WhenSuccess_ExpectNoError() {
	for name, testCase := range s.dbPoolSingleValueTestCases {
		s.T().Run(name, func(t *testing.T) {
			m := &mockConnectionPoolNew{}

			err := testCase.handler(context.Background(), m)

			s.Require().Nil(err)
		})
	}
}

func (s *RepositoryTestSuite) TestGetSingleValue_PropagatesError() {
	for name, testCase := range s.dbPoolSingleValueTestCases {
		s.T().Run(name, func(t *testing.T) {
			m := &mockConnectionPoolNew{
				rows: mockRowsNew{
					getSingleValueErrs: []error{errDefault},
				},
			}

			err := testCase.handler(context.Background(), m)

			s.Require().Equal(errDefault, err)
		})
	}
}

func (s *RepositoryTestSuite) TestGetSingleValue_PropagatesScanError() {
	for name, testCase := range s.dbPoolSingleValueTestCases {
		s.T().Run(name, func(t *testing.T) {
			m := &mockConnectionPoolNew{
				rows: mockRowsNew{
					scanner: &mockScannable{
						err: errDefault,
					},
				},
			}

			err := testCase.handler(context.Background(), m)

			s.Require().Equal(errDefault, err)
		})
	}
}

func (s *RepositoryTestSuite) TestGetSingleValue_ExpectedCorrectPropertiesAreScanned() {
	for name, testCase := range s.dbPoolSingleValueTestCases {
		s.T().Run(name, func(t *testing.T) {
			scanner := &mockScannable{}
			m := &mockConnectionPoolNew{
				rows: mockRowsNew{
					scanner: scanner,
				},
			}

			testCase.handler(context.Background(), m)

			s.Require().Equal(testCase.expectedScanCalls, scanner.scanCalled)
			s.Require().Equal(len(testCase.expectedScannedProps), len(scanner.props))

			for id, expectedProps := range testCase.expectedScannedProps {
				actualProps := scanner.props[id]

				s.Require().Equal(len(expectedProps), len(actualProps))

				for idProp, expectedProp := range expectedProps {
					actualProp := actualProps[idProp]
					s.Require().IsType(expectedProp, actualProp)
				}
			}
		})
	}
}

func (s *RepositoryTestSuite) TestGetAll_ExpectCorrectNumberOfCalls() {
	for name, testCase := range s.dbPoolGetAllTestCases {
		s.T().Run(name, func(t *testing.T) {
			m := &mockConnectionPoolNew{}

			testCase.handler(context.Background(), m)

			s.Require().Equal(testCase.expectedGetAllCalls, m.rows.getAllCalled)
		})
	}
}

func (s *RepositoryTestSuite) TestGetAll_WhenSuccess_ExpectNoError() {
	for name, testCase := range s.dbPoolGetAllTestCases {
		s.T().Run(name, func(t *testing.T) {
			m := &mockConnectionPoolNew{}

			err := testCase.handler(context.Background(), m)

			s.Require().Nil(err)
		})
	}
}

func (s *RepositoryTestSuite) TestGetAll_PropagatesError() {
	for name, testCase := range s.dbPoolGetAllTestCases {
		s.T().Run(name, func(t *testing.T) {
			m := &mockConnectionPoolNew{
				rows: mockRowsNew{
					getAllErrs: []error{errDefault},
				},
			}

			err := testCase.handler(context.Background(), m)

			s.Require().Equal(errDefault, err)
		})
	}
}

func (s *RepositoryTestSuite) TestGetAll_PropagatesScanError() {
	for name, testCase := range s.dbPoolGetAllTestCases {
		s.T().Run(name, func(t *testing.T) {
			m := &mockConnectionPoolNew{
				rows: mockRowsNew{
					scanner: &mockScannable{
						err: errDefault,
					},
				},
			}

			err := testCase.handler(context.Background(), m)

			s.Require().Equal(errDefault, err)
		})
	}
}

func (s *RepositoryTestSuite) TestGetAll_ExpectedCorrectPropertiesAreScanned() {
	for name, testCase := range s.dbPoolGetAllTestCases {
		s.T().Run(name, func(t *testing.T) {
			scanner := &mockScannable{}
			m := &mockConnectionPoolNew{
				rows: mockRowsNew{
					scanner: scanner,
				},
			}

			testCase.handler(context.Background(), m)

			s.Require().Equal(testCase.expectedScanCalls, scanner.scanCalled)
			s.Require().Equal(len(testCase.expectedScannedProps), len(scanner.props))

			for id, expectedProps := range testCase.expectedScannedProps {
				actualProps := scanner.props[id]

				s.Require().Equal(len(expectedProps), len(actualProps))

				for idProp, expectedProp := range expectedProps {
					actualProp := actualProps[idProp]
					s.Require().IsType(expectedProp, actualProp)
				}
			}
		})
	}
}

func (s *RepositoryTestSuite) TestReturnsExpectedValue() {
	for name, testCase := range s.dbPoolReturnTestCases {
		s.T().Run(name, func(t *testing.T) {
			actual := testCase.handler(context.Background(), &mockConnectionPoolNew{})

			s.Require().Equal(testCase.expectedContent, actual)
		})
	}
}

func (tc *dbPoolInteractionTestCase) getCalledCount(m *mockConnectionPoolNew) int {
	switch tc.sqlMode {
	case QueryBased:
		return m.queryCalled
	case ExecBased:
		return m.execCalled
	default:
		panic(fmt.Errorf("Unsupported sql mode %v", tc.sqlMode))
	}
}

func (tc *dbPoolInteractionTestCase) generateErrorMock(err error) *mockConnectionPoolNew {
	switch tc.sqlMode {
	case QueryBased:
		return &mockConnectionPoolNew{
			rows: mockRowsNew{
				err: err,
			},
		}
	case ExecBased:
		return &mockConnectionPoolNew{
			execErr: err,
		}
	default:
		panic(fmt.Errorf("Unsupported sql mode %v", tc.sqlMode))
	}
}
