package service

import (
	"context"
	"testing"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/db"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/repositories"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type testFunc func(context.Context, db.Connection, repositories.Repositories) error
type returnTestFunc func(context.Context, db.Connection, repositories.Repositories) interface{}
type generateRepositoriesMocks func() repositories.Repositories

type verifyError func(error, *require.Assertions)
type verifyMockInteractions func(repositories.Repositories, *require.Assertions)

type generateConnectionMock func() db.Connection
type verifyPoolInteractions func(db.Connection, *require.Assertions)

type repositoryInteractionTestCase struct {
	generateConnectionMock    generateConnectionMock
	generateRepositoriesMocks generateRepositoriesMocks
	handler                   testFunc
	expectedError             error
	verifyError               verifyError
	verifyInteractions        verifyMockInteractions
}

type verifyContent func(interface{}, repositories.Repositories, *require.Assertions)

type returnTestCase struct {
	generateRepositoriesMocks generateRepositoriesMocks
	handler                   returnTestFunc
	expectedContent           interface{}
	verifyContent             verifyContent
}

type transactionTestCase struct {
	generateRepositoriesMocks generateRepositoriesMocks
	handler                   testFunc
}

type transactionInteractionTestCase struct {
	generateConnectionMock    generateConnectionMock
	generateRepositoriesMocks generateRepositoriesMocks
	handler                   testFunc
	expectedError             error
	verifyInteractions        verifyPoolInteractions
	verifyMockInteractions    verifyMockInteractions
}

type ServicePoolTestSuite struct {
	suite.Suite

	generateRepositoriesMocks      generateRepositoriesMocks
	generateErrorRepositoriesMocks generateRepositoriesMocks

	repositoryInteractionTestCases  map[string]repositoryInteractionTestCase
	returnTestCases                 map[string]returnTestCase
	transactionTestCases            map[string]transactionTestCase
	transactionInteractionTestCases map[string]transactionInteractionTestCase
}

func (s *ServicePoolTestSuite) TestWhenCallingHandler_ExpectCorrectInteraction() {
	for name, testCase := range s.repositoryInteractionTestCases {
		s.T().Run(name, func(t *testing.T) {
			var repos repositories.Repositories
			if testCase.generateRepositoriesMocks != nil {
				repos = testCase.generateRepositoriesMocks()
			} else {
				repos = s.generateRepositoriesMocks()
			}

			var conn db.Connection
			if testCase.generateConnectionMock != nil {
				conn = testCase.generateConnectionMock()
			} else {
				conn = &mockConnection{}
			}

			err := testCase.handler(context.Background(), conn, repos)

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
			if testCase.generateRepositoriesMocks != nil {
				repos = testCase.generateRepositoriesMocks()
			} else {
				repos = s.generateRepositoriesMocks()
			}

			actual := testCase.handler(context.Background(), &mockConnection{}, repos)

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
			if testCase.generateRepositoriesMocks != nil {
				repos = testCase.generateRepositoriesMocks()
			} else {
				repos = s.generateRepositoriesMocks()
			}

			m := &mockConnection{}
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
			if testCase.generateRepositoriesMocks != nil {
				repos = testCase.generateRepositoriesMocks()
			} else {
				repos = s.generateRepositoriesMocks()
			}

			m := &mockConnection{
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
			if testCase.generateRepositoriesMocks != nil {
				repos = testCase.generateRepositoriesMocks()
			} else {
				repos = s.generateRepositoriesMocks()
			}

			var m db.Connection
			if testCase.generateConnectionMock != nil {
				m = testCase.generateConnectionMock()
			} else {
				m = &mockConnection{}
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
