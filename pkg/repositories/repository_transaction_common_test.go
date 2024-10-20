package repositories

import (
	"context"
	"testing"

	"github.com/KnoblauchPilze/galactic-sovereign/pkg/db"
	"github.com/stretchr/testify/suite"
)

type dbTransactionInteractionTestCase dbInteractionTestCase[db.Transaction]
type dbTransactionSingleValueTestCase dbSingleValueTestCase[db.Transaction]
type dbTransactionGetAllTestCase dbGetAllTestCase[db.Transaction]
type dbTransactionReturnTestCase dbReturnTestCase[db.Transaction]
type dbTransactionErrorTestCase dbErrorTestCase[db.Transaction]

type RepositoryTransactionTestSuite struct {
	suite.Suite

	dbInteractionTestCases map[string]dbTransactionInteractionTestCase
	dbSingleValueTestCases map[string]dbTransactionSingleValueTestCase
	dbGetAllTestCases      map[string]dbTransactionGetAllTestCase
	dbReturnTestCases      map[string]dbTransactionReturnTestCase
	dbErrorTestCases       map[string]dbTransactionErrorTestCase
}

func (s *RepositoryTransactionTestSuite) TestTransaction_ExpectCorrectNumberOfCalls() {
	for name, testCase := range s.dbInteractionTestCases {
		s.T().Run(name, func(t *testing.T) {
			m := s.generateMockAndAssertType(testCase)

			err := testCase.handler(context.Background(), m)

			s.Require().Nil(err)
			called := getTransactionCalledCount(testCase.sqlMode, m)
			s.Require().Equal(len(testCase.expectedSqlQueries), called)
		})
	}
}

func (s *RepositoryTransactionTestSuite) TestTransaction_ExpectCorrectSqlQueries() {
	for name, testCase := range s.dbInteractionTestCases {
		s.T().Run(name, func(t *testing.T) {
			m := s.generateMockAndAssertType(testCase)

			err := testCase.handler(context.Background(), m)

			s.Require().Nil(err)
			s.Require().Equal(testCase.expectedSqlQueries, m.sqlQueries)
		})
	}
}

func (s *RepositoryTransactionTestSuite) TestTransaction_ExpectCorrectSqlArguments() {
	for name, testCase := range s.dbInteractionTestCases {
		s.T().Run(name, func(t *testing.T) {
			m := s.generateMockAndAssertType(testCase)

			err := testCase.handler(context.Background(), m)

			s.Require().Nil(err)
			s.Require().Equal(len(testCase.expectedArguments), len(m.args))
			for id, expectedArgs := range testCase.expectedArguments {
				actualArgs := m.args[id]
				s.Require().Equal(expectedArgs, actualArgs)
			}
		})
	}
}

func (s *RepositoryTransactionTestSuite) TestTransaction_PropagatesError() {
	for name, testCase := range s.dbInteractionTestCases {
		s.T().Run(name, func(t *testing.T) {
			m := generateTransactionErrorMock(testCase.sqlMode, errDefault)

			err := testCase.handler(context.Background(), m)

			s.Require().Equal(errDefault, err)
		})
	}
}

func (s *RepositoryTransactionTestSuite) TestGetSingleValue_ExpectCorrectNumberOfCalls() {
	for name, testCase := range s.dbSingleValueTestCases {
		s.T().Run(name, func(t *testing.T) {
			m := &mockTransaction{}

			testCase.handler(context.Background(), m)

			s.Require().Equal(testCase.expectedGetSingleValueCalls, m.rows.getSingleValueCalled)
		})
	}
}

func (s *RepositoryTransactionTestSuite) TestGetSingleValue_WhenSuccess_ExpectNoError() {
	for name, testCase := range s.dbSingleValueTestCases {
		s.T().Run(name, func(t *testing.T) {
			m := &mockTransaction{}

			err := testCase.handler(context.Background(), m)

			s.Require().Nil(err)
		})
	}
}

func (s *RepositoryTransactionTestSuite) TestGetSingleValue_PropagatesError() {
	for name, testCase := range s.dbSingleValueTestCases {
		s.T().Run(name, func(t *testing.T) {
			m := &mockTransaction{
				rows: mockRows{
					getSingleValueErrs: []error{errDefault},
				},
			}

			err := testCase.handler(context.Background(), m)

			s.Require().Equal(errDefault, err)
		})
	}
}

func (s *RepositoryTransactionTestSuite) TestGetSingleValue_PropagatesScanError() {
	for name, testCase := range s.dbSingleValueTestCases {
		s.T().Run(name, func(t *testing.T) {
			m := &mockTransaction{
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

func (s *RepositoryTransactionTestSuite) TestGetSingleValue_ExpectedCorrectPropertiesAreScanned() {
	for name, testCase := range s.dbSingleValueTestCases {
		s.T().Run(name, func(t *testing.T) {
			scanner := &mockScannable{}
			m := &mockTransaction{
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

func (s *RepositoryTransactionTestSuite) TestGetAll_ExpectCorrectNumberOfCalls() {
	for name, testCase := range s.dbGetAllTestCases {
		s.T().Run(name, func(t *testing.T) {
			m := &mockTransaction{}

			testCase.handler(context.Background(), m)

			s.Require().Equal(testCase.expectedGetAllCalls, m.rows.getAllCalled)
		})
	}
}

func (s *RepositoryTransactionTestSuite) TestGetAll_WhenSuccess_ExpectNoError() {
	for name, testCase := range s.dbGetAllTestCases {
		s.T().Run(name, func(t *testing.T) {
			m := &mockTransaction{}

			err := testCase.handler(context.Background(), m)

			s.Require().Nil(err)
		})
	}
}

func (s *RepositoryTransactionTestSuite) TestGetAll_PropagatesError() {
	for name, testCase := range s.dbGetAllTestCases {
		s.T().Run(name, func(t *testing.T) {
			m := &mockTransaction{
				rows: mockRows{
					getAllErrs: []error{errDefault},
				},
			}

			err := testCase.handler(context.Background(), m)

			s.Require().Equal(errDefault, err)
		})
	}
}

func (s *RepositoryTransactionTestSuite) TestGetAll_PropagatesScanError() {
	for name, testCase := range s.dbGetAllTestCases {
		s.T().Run(name, func(t *testing.T) {
			m := &mockTransaction{
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

func (s *RepositoryTransactionTestSuite) TestGetAll_ExpectedCorrectPropertiesAreScanned() {
	for name, testCase := range s.dbGetAllTestCases {
		s.T().Run(name, func(t *testing.T) {
			scanner := &mockScannable{}
			m := &mockTransaction{
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

func (s *RepositoryTransactionTestSuite) TestReturnsExpectedValue() {
	for name, testCase := range s.dbReturnTestCases {
		s.T().Run(name, func(t *testing.T) {
			m := &mockTransaction{}

			actual := testCase.handler(context.Background(), m)

			s.Require().Equal(testCase.expectedContent, actual)
		})
	}
}

func (s *RepositoryTransactionTestSuite) TestHandler_ExpectCorrectError() {
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

func (s *RepositoryTransactionTestSuite) generateMockAndAssertType(tc dbTransactionInteractionTestCase) *mockTransaction {
	var out *mockTransaction

	if tc.generateMock == nil {
		out = &mockTransaction{}
	} else {
		maybeMock := tc.generateMock()
		if mock, ok := maybeMock.(*mockTransaction); !ok {
			s.Fail("Transaction mock has not the right type")
		} else {
			out = mock
		}
	}

	return out
}
