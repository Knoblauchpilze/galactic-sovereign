package service

import (
	"context"
	"testing"

	"github.com/KnoblauchPilze/user-service/pkg/db"
	"github.com/KnoblauchPilze/user-service/pkg/repositories"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type testFunc func(context.Context, db.ConnectionPool, repositories.Repositories) error
type returnTestFunc func(context.Context, db.ConnectionPool, repositories.Repositories) interface{}
type generateRepositoriesMock func() repositories.Repositories

type verifyError func(error, *require.Assertions)
type verifyMockInteractions func(repositories.Repositories, *require.Assertions)

type generateConnectionPoolMock func() db.ConnectionPool
type verifyTransactionInteractions func(db.ConnectionPool, *require.Assertions)

type repositoryInteractionTestCase struct {
	generateConnectionPoolMock generateConnectionPoolMock
	generateRepositoriesMock   generateRepositoriesMock
	handler                    testFunc
	expectedError              error
	verifyError                verifyError
	verifyInteractions         verifyMockInteractions
}

type verifyContent func(interface{}, repositories.Repositories, *require.Assertions)

type returnTestCase struct {
	generateRepositoriesMock generateRepositoriesMock
	handler                  returnTestFunc
	expectedContent          interface{}
	verifyContent            verifyContent
}

type transactionTestCase struct {
	generateRepositoriesMock generateRepositoriesMock
	handler                  testFunc
}

type transactionInteractionTestCase struct {
	generateConnectionPoolMock generateConnectionPoolMock
	generateRepositoriesMock   generateRepositoriesMock
	handler                    testFunc
	expectedError              error
	verifyInteractions         verifyTransactionInteractions
	verifyMockInteractions     verifyMockInteractions
}

type ServicePoolTestSuite struct {
	suite.Suite

	generateRepositoriesMock      generateRepositoriesMock
	generateErrorRepositoriesMock generateRepositoriesMock

	repositoryInteractionTestCases  map[string]repositoryInteractionTestCase
	returnTestCases                 map[string]returnTestCase
	transactionTestCases            map[string]transactionTestCase
	transactionInteractionTestCases map[string]transactionInteractionTestCase
}

func (s *ServicePoolTestSuite) TestWhenCallingHandler_ExpectCorrectInteraction() {
	for name, testCase := range s.repositoryInteractionTestCases {
		s.T().Run(name, func(t *testing.T) {
			var repos repositories.Repositories
			if testCase.generateRepositoriesMock != nil {
				repos = testCase.generateRepositoriesMock()
			} else {
				repos = s.generateRepositoriesMock()
			}

			var pool db.ConnectionPool
			if testCase.generateConnectionPoolMock != nil {
				pool = testCase.generateConnectionPoolMock()
			} else {
				pool = &mockConnectionPool{}
			}

			err := testCase.handler(context.Background(), pool, repos)

			if testCase.verifyError != nil {
				testCase.verifyError(err, s.Require())
			} else {
				s.Require().Equal(testCase.expectedError, err)
			}
			if testCase.verifyInteractions != nil {
				testCase.verifyInteractions(repos, s.Require())
			}
		})
	}
}

func (s *ServicePoolTestSuite) TestWhenRepositorySucceeds_ReturnsExpectedValue() {
	for name, testCase := range s.returnTestCases {
		s.T().Run(name, func(t *testing.T) {
			var repos repositories.Repositories
			if testCase.generateRepositoriesMock != nil {
				repos = testCase.generateRepositoriesMock()
			} else {
				repos = s.generateRepositoriesMock()
			}

			actual := testCase.handler(context.Background(), &mockConnectionPool{}, repos)

			if testCase.verifyContent != nil {
				testCase.verifyContent(actual, repos, s.Require())
			} else {
				s.Require().Equal(testCase.expectedContent, actual)
			}
		})
	}
}

func (s *ServicePoolTestSuite) TestWhenUsingTransaction_ExpectCallsClose() {
	for name, testCase := range s.transactionTestCases {
		s.T().Run(name, func(t *testing.T) {
			var repos repositories.Repositories
			if testCase.generateRepositoriesMock != nil {
				repos = testCase.generateRepositoriesMock()
			} else {
				repos = s.generateRepositoriesMock()
			}

			m := &mockConnectionPool{}
			testCase.handler(context.Background(), m, repos)

			for _, tx := range m.txs {
				s.Require().Equal(1, tx.closeCalled)
			}
		})
	}
}

func (s *ServicePoolTestSuite) TestWhenCreatingTransactionFails_ExpectErrorIsPropagated() {
	for name, testCase := range s.transactionTestCases {
		s.T().Run(name, func(t *testing.T) {
			var repos repositories.Repositories
			if testCase.generateRepositoriesMock != nil {
				repos = testCase.generateRepositoriesMock()
			} else {
				repos = s.generateRepositoriesMock()
			}

			m := &mockConnectionPool{
				errs: []error{errDefault},
			}
			err := testCase.handler(context.Background(), m, repos)

			s.Require().Equal(errDefault, err)
		})
	}
}

func (s *ServicePoolTestSuite) TestWhenUsingTransaction_ExpectCorrectInteraction() {
	for name, testCase := range s.transactionInteractionTestCases {
		s.T().Run(name, func(t *testing.T) {
			var repos repositories.Repositories
			if testCase.generateRepositoriesMock != nil {
				repos = testCase.generateRepositoriesMock()
			} else {
				repos = s.generateRepositoriesMock()
			}

			var m db.ConnectionPool
			if testCase.generateConnectionPoolMock != nil {
				m = testCase.generateConnectionPoolMock()
			} else {
				m = &mockConnectionPool{}
			}
			err := testCase.handler(context.Background(), m, repos)

			s.Require().Equal(testCase.expectedError, err)
			if testCase.verifyInteractions != nil {
				testCase.verifyInteractions(m, s.Require())
			}
			if testCase.verifyMockInteractions != nil {
				testCase.verifyMockInteractions(repos, s.Require())
			}
		})
	}
}
