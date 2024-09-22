package service

import (
	"context"
	"testing"

	"github.com/KnoblauchPilze/user-service/pkg/db"
	"github.com/KnoblauchPilze/user-service/pkg/repositories"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type testTransactionFunc func(context.Context, db.Transaction, repositories.Repositories) error

type generateTransactionMock func() db.Transaction
type verifyTransactionInteractions func(db.Transaction, *require.Assertions)

type serviceTransactionInteractionTestCase struct {
	generateTransactionMock       generateTransactionMock
	generateRepositoriesMock      generateRepositoriesMock
	handler                       testTransactionFunc
	expectedError                 error
	verifyError                   verifyError
	verifyInteractions            verifyMockInteractions
	verifyTransactionInteractions verifyTransactionInteractions
}

type ServiceTransactionTestSuite struct {
	suite.Suite

	generateRepositoriesMock generateRepositoriesMock

	interactionTestCases map[string]serviceTransactionInteractionTestCase
}

func (s *ServiceTransactionTestSuite) TestWhenCallingHandler_ExpectCorrectInteraction() {
	for name, testCase := range s.interactionTestCases {
		s.T().Run(name, func(t *testing.T) {
			var repos repositories.Repositories
			if testCase.generateRepositoriesMock != nil {
				repos = testCase.generateRepositoriesMock()
			} else {
				repos = s.generateRepositoriesMock()
			}

			var tx db.Transaction
			if testCase.generateTransactionMock != nil {
				tx = testCase.generateTransactionMock()
			} else {
				tx = &mockTransaction{}
			}

			err := testCase.handler(context.Background(), tx, repos)

			if testCase.verifyError != nil {
				testCase.verifyError(err, s.Require())
			} else {
				s.Require().Equal(testCase.expectedError, err)
			}
			if testCase.verifyInteractions != nil {
				testCase.verifyInteractions(repos, s.Require())
			}
			if testCase.verifyTransactionInteractions != nil {
				testCase.verifyTransactionInteractions(tx, s.Require())
			}
		})
	}
}
