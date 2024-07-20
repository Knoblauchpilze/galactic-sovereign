package repositories

import (
	"context"
	"testing"

	"github.com/KnoblauchPilze/user-service/pkg/db"
	"github.com/stretchr/testify/suite"
)

type dbTransactionInteractionTestCase dbInteractionTestCase[db.Transaction]
type dbTransactionSingleValueTestCase dbSingleValueTestCase[db.Transaction]
type dbTransactionGetAllTestCase dbGetAllTestCase[db.Transaction]
type dbTransactionReturnTestCase dbReturnTestCase[db.Transaction]
type dbTransactionErrorTestCase dbErrorTestCase[db.Transaction]

type RepositoryTransactionTestSuiteNew struct {
	suite.Suite

	dbInteractionTestCases map[string]dbTransactionInteractionTestCase
	dbSingleValueTestCases map[string]dbTransactionSingleValueTestCase
	dbGetAllTestCases      map[string]dbTransactionGetAllTestCase
	dbReturnTestCases      map[string]dbTransactionReturnTestCase
	dbErrorTestCases       map[string]dbTransactionErrorTestCase
}

func (s *RepositoryTransactionTestSuiteNew) TestTransaction_ExpectCorrectNumberOfCalls() {
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

func (s *RepositoryTransactionTestSuiteNew) TestTransaction_ExpectCorrectSqlQueries() {
	for name, testCase := range s.dbInteractionTestCases {
		s.T().Run(name, func(t *testing.T) {
			m := s.generateMockAndAssertType(testCase)

			err := testCase.handler(context.Background(), m)

			s.Require().Nil(err)
			s.Require().Equal(testCase.expectedSqlQueries, m.sqlQueries)
		})
	}
}

func (s *RepositoryTransactionTestSuiteNew) TestTransaction_ExpectCorrectSqlArguments() {
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

func (s *RepositoryTransactionTestSuiteNew) TestTransaction_PropagatesError() {
	for name, testCase := range s.dbInteractionTestCases {
		s.T().Run(name, func(t *testing.T) {
			m := generateTransactionErrorMock(testCase.sqlMode, errDefault)

			err := testCase.handler(context.Background(), m)

			s.Require().Equal(errDefault, err)
		})
	}
}

func (s *RepositoryTransactionTestSuiteNew) TestGetSingleValue_ExpectCorrectNumberOfCalls() {
	for name, testCase := range s.dbSingleValueTestCases {
		s.T().Run(name, func(t *testing.T) {
			m := &mockTransactionNew{}

			testCase.handler(context.Background(), m)

			s.Require().Equal(testCase.expectedGetSingleValueCalls, m.rows.getSingleValueCalled)
		})
	}
}

func (s *RepositoryTransactionTestSuiteNew) TestGetSingleValue_WhenSuccess_ExpectNoError() {
	for name, testCase := range s.dbSingleValueTestCases {
		s.T().Run(name, func(t *testing.T) {
			m := &mockTransactionNew{}

			err := testCase.handler(context.Background(), m)

			s.Require().Nil(err)
		})
	}
}

func (s *RepositoryTransactionTestSuiteNew) TestGetSingleValue_PropagatesError() {
	for name, testCase := range s.dbSingleValueTestCases {
		s.T().Run(name, func(t *testing.T) {
			m := &mockTransactionNew{
				rows: mockRowsNew{
					getSingleValueErrs: []error{errDefault},
				},
			}

			err := testCase.handler(context.Background(), m)

			s.Require().Equal(errDefault, err)
		})
	}
}

func (s *RepositoryTransactionTestSuiteNew) TestGetSingleValue_PropagatesScanError() {
	for name, testCase := range s.dbSingleValueTestCases {
		s.T().Run(name, func(t *testing.T) {
			m := &mockTransactionNew{
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

func (s *RepositoryTransactionTestSuiteNew) TestGetSingleValue_ExpectedCorrectPropertiesAreScanned() {
	for name, testCase := range s.dbSingleValueTestCases {
		s.T().Run(name, func(t *testing.T) {
			scanner := &mockScannable{}
			m := &mockTransactionNew{
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

func (s *RepositoryTransactionTestSuiteNew) TestGetAll_ExpectCorrectNumberOfCalls() {
	for name, testCase := range s.dbGetAllTestCases {
		s.T().Run(name, func(t *testing.T) {
			m := &mockTransactionNew{}

			testCase.handler(context.Background(), m)

			s.Require().Equal(testCase.expectedGetAllCalls, m.rows.getAllCalled)
		})
	}
}

func (s *RepositoryTransactionTestSuiteNew) TestGetAll_WhenSuccess_ExpectNoError() {
	for name, testCase := range s.dbGetAllTestCases {
		s.T().Run(name, func(t *testing.T) {
			m := &mockTransactionNew{}

			err := testCase.handler(context.Background(), m)

			s.Require().Nil(err)
		})
	}
}

func (s *RepositoryTransactionTestSuiteNew) TestGetAll_PropagatesError() {
	for name, testCase := range s.dbGetAllTestCases {
		s.T().Run(name, func(t *testing.T) {
			m := &mockTransactionNew{
				rows: mockRowsNew{
					getAllErrs: []error{errDefault},
				},
			}

			err := testCase.handler(context.Background(), m)

			s.Require().Equal(errDefault, err)
		})
	}
}

func (s *RepositoryTransactionTestSuiteNew) TestGetAll_PropagatesScanError() {
	for name, testCase := range s.dbGetAllTestCases {
		s.T().Run(name, func(t *testing.T) {
			m := &mockTransactionNew{
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

func (s *RepositoryTransactionTestSuiteNew) TestGetAll_ExpectedCorrectPropertiesAreScanned() {
	for name, testCase := range s.dbGetAllTestCases {
		s.T().Run(name, func(t *testing.T) {
			scanner := &mockScannable{}
			m := &mockTransactionNew{
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

func (s *RepositoryTransactionTestSuiteNew) TestReturnsExpectedValue() {
	for name, testCase := range s.dbReturnTestCases {
		s.T().Run(name, func(t *testing.T) {
			m := &mockTransactionNew{}

			actual := testCase.handler(context.Background(), m)

			s.Require().Equal(testCase.expectedContent, actual)
		})
	}
}

func (s *RepositoryTransactionTestSuiteNew) TestHandler_ExpectCorrectError() {
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

func (s *RepositoryTransactionTestSuiteNew) generateMockAndAssertType(tc dbTransactionInteractionTestCase) *mockTransactionNew {
	var out *mockTransactionNew

	if tc.generateMock == nil {
		out = &mockTransactionNew{}
	} else {
		maybeMock := tc.generateMock()
		if mock, ok := maybeMock.(*mockTransactionNew); !ok {
			s.Fail("Transaction mock has not the right type")
		} else {
			out = mock
		}
	}

	return out
}
