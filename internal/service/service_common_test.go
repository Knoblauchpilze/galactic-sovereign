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
type generateValidRepositoriesMock func() repositories.Repositories
type generateErrorRepositoriesMock func(err error) repositories.Repositories

type verifyError func(error, *require.Assertions)

type errorTestCase struct {
	generateErrorRepositoriesMock generateErrorRepositoriesMock
	handler                       testFunc
	verifyError                   verifyError
}

type verifyMockInteractions func(repositories.Repositories, *require.Assertions)

type repositoryInteractionTestCase struct {
	generateValidRepositoriesMock generateValidRepositoriesMock
	handler                       testFunc
	expectedError                 error
	verifyInteractions            verifyMockInteractions
}

type verifyContent func(interface{}, repositories.Repositories, *require.Assertions)

type returnTestCase struct {
	generateValidRepositoriesMock generateValidRepositoriesMock
	handler                       returnTestFunc
	expectedContent               interface{}
	verifyContent                 verifyContent
}

type transactionTestCase struct {
	generateValidRepositoriesMock generateValidRepositoriesMock
	handler                       testFunc
}

type ServiceTestSuite struct {
	suite.Suite

	generateValidRepositoriesMock generateValidRepositoriesMock
	generateErrorRepositoriesMock generateErrorRepositoriesMock

	errorTestCases                 map[string]errorTestCase
	repositoryInteractionTestCases map[string]repositoryInteractionTestCase
	returnTestCases                map[string]returnTestCase
	transactionTestCases           map[string]transactionTestCase
}

func (s *ServiceTestSuite) TestWhenRepositoryFails_ExpectErrorIsPropagated() {
	for name, testCase := range s.errorTestCases {
		s.T().Run(name, func(t *testing.T) {
			var repos repositories.Repositories
			if testCase.generateErrorRepositoriesMock != nil {
				repos = testCase.generateErrorRepositoriesMock(errDefault)
			} else {
				repos = s.generateErrorRepositoriesMock(errDefault)
			}

			err := testCase.handler(context.Background(), &mockConnectionPool{}, repos)

			if testCase.verifyError != nil {
				testCase.verifyError(err, s.Require())
			} else {
				s.Require().Equal(errDefault, err)
			}
		})
	}
}

func (s *ServiceTestSuite) TestWhenRepositorySucceeds_ExpectCorrectInteraction() {
	for name, testCase := range s.repositoryInteractionTestCases {
		s.T().Run(name, func(t *testing.T) {
			var repos repositories.Repositories
			if testCase.generateValidRepositoriesMock != nil {
				repos = testCase.generateValidRepositoriesMock()
			} else {
				repos = s.generateValidRepositoriesMock()
			}

			err := testCase.handler(context.Background(), &mockConnectionPool{}, repos)

			s.Require().Equal(testCase.expectedError, err)
			testCase.verifyInteractions(repos, s.Require())
		})
	}
}

func (s *ServiceTestSuite) TestWhenRepositorySucceeds_ReturnsExpectedValue() {
	for name, testCase := range s.returnTestCases {
		s.T().Run(name, func(t *testing.T) {
			var repos repositories.Repositories
			if testCase.generateValidRepositoriesMock != nil {
				repos = testCase.generateValidRepositoriesMock()
			} else {
				repos = s.generateValidRepositoriesMock()
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

func (s *ServiceTestSuite) TestWhenUsingTransaction_ExpectCallsClose() {
	for name, testCase := range s.transactionTestCases {
		s.T().Run(name, func(t *testing.T) {
			var repos repositories.Repositories
			if testCase.generateValidRepositoriesMock != nil {
				repos = testCase.generateValidRepositoriesMock()
			} else {
				repos = s.generateValidRepositoriesMock()
			}

			m := &mockConnectionPool{}
			testCase.handler(context.Background(), m, repos)

			s.Require().Equal(1, m.tx.closeCalled)
		})
	}
}

func (s *ServiceTestSuite) TestWhenCreatingTransactionFails_ExpectErrorIsPropagated() {
	for name, testCase := range s.transactionTestCases {
		s.T().Run(name, func(t *testing.T) {
			var repos repositories.Repositories
			if testCase.generateValidRepositoriesMock != nil {
				repos = testCase.generateValidRepositoriesMock()
			} else {
				repos = s.generateValidRepositoriesMock()
			}

			m := &mockConnectionPool{
				err: errDefault,
			}
			err := testCase.handler(context.Background(), m, repos)

			s.Require().Equal(errDefault, err)
		})
	}
}
