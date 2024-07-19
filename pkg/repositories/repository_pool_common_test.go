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
			m := &mockConnectionPoolNew{}

			testCase.handler(context.Background(), m)

			called := getPoolCalledCount(testCase.sqlMode, m)
			s.Require().Equal(1, called)
		})
	}
}

func (s *RepositoryPoolTestSuite) TestPool_ExpectCorrectSqlQuery() {
	for name, testCase := range s.dbInteractionTestCases {
		s.T().Run(name, func(t *testing.T) {
			m := &mockConnectionPoolNew{}

			testCase.handler(context.Background(), m)

			s.Require().Equal(testCase.expectedSql, m.sqlQuery)
		})
	}
}

func (s *RepositoryPoolTestSuite) TestPool_ExpectCorrectSqlArguments() {
	for name, testCase := range s.dbInteractionTestCases {
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
