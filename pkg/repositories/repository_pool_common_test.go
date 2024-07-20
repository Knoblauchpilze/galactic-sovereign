package repositories

import (
	"context"
	"testing"

	"github.com/KnoblauchPilze/user-service/pkg/db"
	"github.com/stretchr/testify/suite"
)

type dbPoolInteractionTestCase dbInteractionTestCase[db.ConnectionPool]
type dbPoolSingleValueTestCase dbSingleValueTestCase[db.ConnectionPool]
type dbPoolGetAllTestCase dbGetAllTestCase[db.ConnectionPool]
type dbPoolReturnTestCase dbReturnTestCase[db.ConnectionPool]

type RepositoryPoolTestSuite struct {
	suite.Suite

	dbInteractionTestCases map[string]dbPoolInteractionTestCase
	dbSingleValueTestCases map[string]dbPoolSingleValueTestCase
	dbGetAllTestCases      map[string]dbPoolGetAllTestCase
	dbReturnTestCases      map[string]dbPoolReturnTestCase
}

func (s *RepositoryPoolTestSuite) TestPool_ExpectCorrectNumberOfCalls() {
	for name, testCase := range s.dbInteractionTestCases {
		s.T().Run(name, func(t *testing.T) {
			m := s.generateMockAndAssertType(testCase)

			err := testCase.handler(context.Background(), m)

			s.Require().Nil(err)
			called := getPoolCalledCount(testCase.sqlMode, m)
			s.Require().Equal(1, called)
		})
	}
}

func (s *RepositoryPoolTestSuite) TestPool_ExpectCorrectSqlQuery() {
	for name, testCase := range s.dbInteractionTestCases {
		s.T().Run(name, func(t *testing.T) {
			m := s.generateMockAndAssertType(testCase)

			err := testCase.handler(context.Background(), m)

			s.Require().Nil(err)
			s.Require().Equal(1, len(testCase.expectedSqlQueries))
			s.Require().Equal(testCase.expectedSqlQueries[0], m.sqlQuery)
		})
	}
}

func (s *RepositoryPoolTestSuite) TestPool_ExpectCorrectSqlArguments() {
	for name, testCase := range s.dbInteractionTestCases {
		s.T().Run(name, func(t *testing.T) {
			m := s.generateMockAndAssertType(testCase)

			err := testCase.handler(context.Background(), m)

			s.Require().Nil(err)
			s.Require().GreaterOrEqual(1, len(testCase.expectedArguments))
			if len(testCase.expectedArguments) == 0 {
				s.Require().Nil(m.args)
			} else {
				expectedArguments := testCase.expectedArguments[0]
				s.Require().Equal(len(expectedArguments), len(m.args))
				for id, expected := range expectedArguments {
					actual := m.args[id]
					s.Require().Equal(expected, actual)
				}
			}
		})
	}
}

func (s *RepositoryPoolTestSuite) TestPool_PropagatesError() {
	for name, testCase := range s.dbInteractionTestCases {
		s.T().Run(name, func(t *testing.T) {
			m := generatePoolErrorMock(testCase.sqlMode, errDefault)

			err := testCase.handler(context.Background(), m)

			s.Require().Equal(errDefault, err)
		})
	}
}

func (s *RepositoryPoolTestSuite) TestGetSingleValue_ExpectCorrectNumberOfCalls() {
	for name, testCase := range s.dbSingleValueTestCases {
		s.T().Run(name, func(t *testing.T) {
			m := &mockConnectionPoolNew{}

			testCase.handler(context.Background(), m)

			s.Require().Equal(testCase.expectedGetSingleValueCalls, m.rows.getSingleValueCalled)
		})
	}
}

func (s *RepositoryPoolTestSuite) TestGetSingleValue_WhenSuccess_ExpectNoError() {
	for name, testCase := range s.dbSingleValueTestCases {
		s.T().Run(name, func(t *testing.T) {
			m := &mockConnectionPoolNew{}

			err := testCase.handler(context.Background(), m)

			s.Require().Nil(err)
		})
	}
}

func (s *RepositoryPoolTestSuite) TestGetSingleValue_PropagatesError() {
	for name, testCase := range s.dbSingleValueTestCases {
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

func (s *RepositoryPoolTestSuite) TestGetSingleValue_PropagatesScanError() {
	for name, testCase := range s.dbSingleValueTestCases {
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

func (s *RepositoryPoolTestSuite) TestGetSingleValue_ExpectedCorrectPropertiesAreScanned() {
	for name, testCase := range s.dbSingleValueTestCases {
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

func (s *RepositoryPoolTestSuite) TestGetAll_ExpectCorrectNumberOfCalls() {
	for name, testCase := range s.dbGetAllTestCases {
		s.T().Run(name, func(t *testing.T) {
			m := &mockConnectionPoolNew{}

			testCase.handler(context.Background(), m)

			s.Require().Equal(testCase.expectedGetAllCalls, m.rows.getAllCalled)
		})
	}
}

func (s *RepositoryPoolTestSuite) TestGetAll_WhenSuccess_ExpectNoError() {
	for name, testCase := range s.dbGetAllTestCases {
		s.T().Run(name, func(t *testing.T) {
			m := &mockConnectionPoolNew{}

			err := testCase.handler(context.Background(), m)

			s.Require().Nil(err)
		})
	}
}

func (s *RepositoryPoolTestSuite) TestGetAll_PropagatesError() {
	for name, testCase := range s.dbGetAllTestCases {
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

func (s *RepositoryPoolTestSuite) TestGetAll_PropagatesScanError() {
	for name, testCase := range s.dbGetAllTestCases {
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

func (s *RepositoryPoolTestSuite) TestGetAll_ExpectedCorrectPropertiesAreScanned() {
	for name, testCase := range s.dbGetAllTestCases {
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

func (s *RepositoryPoolTestSuite) TestReturnsExpectedValue() {
	for name, testCase := range s.dbReturnTestCases {
		s.T().Run(name, func(t *testing.T) {
			actual := testCase.handler(context.Background(), &mockConnectionPoolNew{})

			s.Require().Equal(testCase.expectedContent, actual)
		})
	}
}

func (s *RepositoryPoolTestSuite) generateMockAndAssertType(tc dbPoolInteractionTestCase) *mockConnectionPoolNew {
	var out *mockConnectionPoolNew

	if tc.generateMock == nil {
		out = &mockConnectionPoolNew{}
	} else {
		maybeMock := tc.generateMock()
		if mock, ok := maybeMock.(*mockConnectionPoolNew); !ok {
			s.Fail("Connection pool mock has not the right type")
		} else {
			out = mock
		}
	}

	return out
}
