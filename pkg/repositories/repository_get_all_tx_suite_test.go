package repositories

import (
	"context"

	"github.com/KnoblauchPilze/user-service/pkg/db"
	"github.com/stretchr/testify/suite"
)

type testGetAllTxFunc func(context.Context, db.Transaction) error

type RepositoryGetAllTransactionTestSuite struct {
	suite.Suite

	testFunc testGetAllTxFunc

	expectedScanCalls    int
	expectedScannedProps [][]interface{}
}

func (s *RepositoryGetAllTransactionTestSuite) TestCallsGetAll() {
	mock := &mockTransaction{}

	s.testFunc(context.Background(), mock)

	s.Require().Equal(1, mock.rows.allCalled)
}

func (s *RepositoryGetAllTransactionTestSuite) TestPropagatesGetAllError() {
	mock := &mockTransaction{
		rows: mockRows{
			allErr: errDefault,
		},
	}

	err := s.testFunc(context.Background(), mock)

	s.Require().Equal(errDefault, err)
}

func (s *RepositoryGetAllTransactionTestSuite) TestPropagatesScanError() {
	mock := &mockTransaction{
		rows: mockRows{
			scanner: &mockScannable{
				err: errDefault,
			},
		},
	}

	err := s.testFunc(context.Background(), mock)

	s.Require().Equal(errDefault, err)
}

func (s *RepositoryGetAllTransactionTestSuite) TestWhenGetAllSucceedsExpectsNoError() {
	mock := &mockTransaction{}

	err := s.testFunc(context.Background(), mock)

	s.Require().Nil(err)
}

func (s *RepositoryGetAllTransactionTestSuite) TestScansExpectedProperties() {
	mock := &mockTransaction{
		rows: mockRows{
			scanner: &mockScannable{},
		},
	}

	s.testFunc(context.Background(), mock)

	s.Require().Equal(s.expectedScanCalls, mock.rows.scanner.scanCalled)
	s.Require().Equal(len(s.expectedScannedProps), len(mock.rows.scanner.props))

	for id, expectedProps := range s.expectedScannedProps {
		actualProps := mock.rows.scanner.props[id]

		s.Require().Equal(len(expectedProps), len(actualProps))

		for idProp, expectedProp := range expectedProps {
			actualProp := actualProps[idProp]
			s.Require().IsType(expectedProp, actualProp)
		}
	}
}
