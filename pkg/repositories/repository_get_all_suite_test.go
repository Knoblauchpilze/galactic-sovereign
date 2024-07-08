package repositories

import (
	"context"

	"github.com/KnoblauchPilze/user-service/pkg/db"
	"github.com/stretchr/testify/assert"
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
	assert := assert.New(s.T())

	mock := &mockConnectionPool{}

	s.testFunc(context.Background(), mock)

	assert.Equal(1, mock.rows.allCalled)
}

func (s *RepositoryGetAllTestSuite) TestPropagatesGetAllError() {
	assert := assert.New(s.T())

	mock := &mockConnectionPool{
		rows: mockRows{
			allErr: errDefault,
		},
	}

	err := s.testFunc(context.Background(), mock)

	assert.Equal(errDefault, err)
}

func (s *RepositoryGetAllTestSuite) TestPropagatesScanError() {
	assert := assert.New(s.T())

	mock := &mockConnectionPool{
		rows: mockRows{
			scanner: &mockScannable{
				err: errDefault,
			},
		},
	}

	err := s.testFunc(context.Background(), mock)

	assert.Equal(errDefault, err)
}

func (s *RepositoryGetAllTestSuite) TestWhenGetAllSucceedsExpectsNoError() {
	assert := assert.New(s.T())

	mock := &mockConnectionPool{}

	err := s.testFunc(context.Background(), mock)

	assert.Nil(err)
}

func (s *RepositoryGetAllTestSuite) TestScansExpectedProperties() {
	assert := assert.New(s.T())

	mock := &mockConnectionPool{
		rows: mockRows{
			scanner: &mockScannable{},
		},
	}

	s.testFunc(context.Background(), mock)

	assert.Equal(s.expectedScanCalls, mock.rows.scanner.scanCalled)
	assert.Equal(len(s.expectedScannedProps), len(mock.rows.scanner.props))

	for id, expectedProps := range s.expectedScannedProps {
		actualProps := mock.rows.scanner.props[id]

		assert.Equal(len(expectedProps), len(actualProps))

		for idProp, expectedProp := range expectedProps {
			actualProp := actualProps[idProp]
			assert.IsType(expectedProp, actualProp)
		}
	}
}
