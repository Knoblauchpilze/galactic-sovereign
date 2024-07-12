package repositories

import (
	"context"

	"github.com/KnoblauchPilze/user-service/pkg/db"
	"github.com/stretchr/testify/suite"
)

type testGetAllFunc func(context.Context, db.ConnectionPool) error

type RepositoryGetAllTestSuite struct {
	suite.Suite

	testFunc testGetAllFunc

	expectedScanCalls    int
	expectedScannedProps [][]interface{}
}

func (s *RepositoryGetAllTestSuite) TestCallsGetAll() {
	mock := &mockConnectionPool{}

	s.testFunc(context.Background(), mock)

	s.Require().Equal(1, mock.rows.allCalled)
}

func (s *RepositoryGetAllTestSuite) TestPropagatesGetAllError() {
	mock := &mockConnectionPool{
		rows: mockRows{
			allErr: errDefault,
		},
	}

	err := s.testFunc(context.Background(), mock)

	s.Require().Equal(errDefault, err)
}

func (s *RepositoryGetAllTestSuite) TestPropagatesScanError() {
	mock := &mockConnectionPool{
		rows: mockRows{
			scanner: &mockScannable{
				err: errDefault,
			},
		},
	}

	err := s.testFunc(context.Background(), mock)

	s.Require().Equal(errDefault, err)
}

func (s *RepositoryGetAllTestSuite) TestWhenGetAllSucceedsExpectsNoError() {
	mock := &mockConnectionPool{}

	err := s.testFunc(context.Background(), mock)

	s.Require().Nil(err)
}

func (s *RepositoryGetAllTestSuite) TestScansExpectedProperties() {
	mock := &mockConnectionPool{
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
