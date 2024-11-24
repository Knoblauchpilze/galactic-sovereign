package repositories

import (
	"context"
	"fmt"
	"testing"

	"github.com/KnoblauchPilze/galactic-sovereign/pkg/db"
	"github.com/stretchr/testify/suite"
)

type dbPoolInteractionTestCase dbInteractionTestCase[db.ConnectionPool]
type dbPoolSingleValueTestCase dbSingleValueTestCase[db.ConnectionPool]
type dbPoolGetAllTestCase dbGetAllTestCase[db.ConnectionPool]
type dbPoolReturnTestCase dbReturnTestCase[db.ConnectionPool]
type dbPoolErrorTestCase dbErrorTestCase[db.ConnectionPool]

type RepositoryPoolTestSuite struct {
	suite.Suite

	dbInteractionTestCases map[string]dbPoolInteractionTestCase
	dbSingleValueTestCases map[string]dbPoolSingleValueTestCase
	dbGetAllTestCases      map[string]dbPoolGetAllTestCase
	dbReturnTestCases      map[string]dbPoolReturnTestCase
	dbErrorTestCases       map[string]dbPoolErrorTestCase
}

var errDefault = fmt.Errorf("some error")

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
			m := &mockConnectionPool{}

			testCase.handler(context.Background(), m)

			s.Require().Equal(testCase.expectedGetSingleValueCalls, m.rows.getSingleValueCalled)
		})
	}
}

func (s *RepositoryPoolTestSuite) TestGetSingleValue_WhenSuccess_ExpectNoError() {
	for name, testCase := range s.dbSingleValueTestCases {
		s.T().Run(name, func(t *testing.T) {
			m := &mockConnectionPool{}

			err := testCase.handler(context.Background(), m)

			s.Require().Nil(err)
		})
	}
}

func (s *RepositoryPoolTestSuite) TestGetSingleValue_PropagatesError() {
	for name, testCase := range s.dbSingleValueTestCases {
		s.T().Run(name, func(t *testing.T) {
			m := &mockConnectionPool{
				rows: mockRows{
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
			m := &mockConnectionPool{
				rows: mockRows{
					scanner: &mockScannable{
						errs: []error{errDefault},
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
			m := &mockConnectionPool{
				rows: mockRows{
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
			m := &mockConnectionPool{}

			testCase.handler(context.Background(), m)

			s.Require().Equal(testCase.expectedGetAllCalls, m.rows.getAllCalled)
		})
	}
}

func (s *RepositoryPoolTestSuite) TestGetAll_WhenSuccess_ExpectNoError() {
	for name, testCase := range s.dbGetAllTestCases {
		s.T().Run(name, func(t *testing.T) {
			m := &mockConnectionPool{}

			err := testCase.handler(context.Background(), m)

			s.Require().Nil(err)
		})
	}
}

func (s *RepositoryPoolTestSuite) TestGetAll_PropagatesError() {
	for name, testCase := range s.dbGetAllTestCases {
		s.T().Run(name, func(t *testing.T) {
			m := &mockConnectionPool{
				rows: mockRows{
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
			m := &mockConnectionPool{
				rows: mockRows{
					scanner: &mockScannable{
						errs: []error{errDefault},
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
			m := &mockConnectionPool{
				rows: mockRows{
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
			actual := testCase.handler(context.Background(), &mockConnectionPool{})

			s.Require().Equal(testCase.expectedContent, actual)
		})
	}
}

func (s *RepositoryPoolTestSuite) TestHandler_ExpectCorrectError() {
	for name, testCase := range s.dbErrorTestCases {
		s.T().Run(name, func(t *testing.T) {
			m := testCase.generateMock()

			err := testCase.handler(context.Background(), m)

			if testCase.verifyError == nil {
				s.Require().Equal(testCase.expectedError, err)
			} else {
				testCase.verifyError(err, s.Require())
			}
		})
	}
}

func (s *RepositoryPoolTestSuite) generateMockAndAssertType(tc dbPoolInteractionTestCase) *mockConnectionPool {
	var out *mockConnectionPool

	if tc.generateMock == nil {
		out = &mockConnectionPool{}
	} else {
		maybeMock := tc.generateMock()
		if mock, ok := maybeMock.(*mockConnectionPool); !ok {
			s.Fail("Connection pool mock has not the right type")
		} else {
			out = mock
		}
	}

	return out
}
