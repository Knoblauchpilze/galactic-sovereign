package repositories

import (
	"context"

	"github.com/KnoblauchPilze/user-service/pkg/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type testGetAllTxFunc func(context.Context, db.Transaction) error

type RepositoryGetAllTransactionTestSuite struct {
	suite.Suite

	testFunc testGetAllTxFunc

	expectedScanCalls    int
	expectedScannedProps []interface{}
}

func (s *RepositoryGetAllTransactionTestSuite) TestCallsGetSingleValue() {
	assert := assert.New(s.T())

	mock := &mockTransaction{}

	s.testFunc(context.Background(), mock)

	assert.Equal(1, mock.rows.allCalled)
}

func (s *RepositoryGetAllTransactionTestSuite) TestPropagatesGetAllError() {
	assert := assert.New(s.T())

	mock := &mockTransaction{
		rows: mockRows{
			allErr: errDefault,
		},
	}

	err := s.testFunc(context.Background(), mock)

	assert.Equal(errDefault, err)
}

func (s *RepositoryGetAllTransactionTestSuite) TestPropagatesScanError() {
	assert := assert.New(s.T())

	mock := &mockTransaction{
		rows: mockRows{
			scanner: &mockScannable{
				err: errDefault,
			},
		},
	}

	err := s.testFunc(context.Background(), mock)

	assert.Equal(errDefault, err)
}

func (s *RepositoryGetAllTransactionTestSuite) TestWhenSingleValueSucceedsExpectsNoError() {
	assert := assert.New(s.T())

	mock := &mockTransaction{}

	err := s.testFunc(context.Background(), mock)

	assert.Nil(err)
}

func (s *RepositoryGetAllTransactionTestSuite) TestScansExpectedProperties() {
	assert := assert.New(s.T())

	mock := &mockTransaction{
		rows: mockRows{
			scanner: &mockScannable{},
		},
	}

	s.testFunc(context.Background(), mock)

	assert.Equal(s.expectedScanCalls, mock.rows.scanner.scanCalled)
	assert.Equal(len(s.expectedScannedProps), len(mock.rows.scanner.props))
	for id, expected := range s.expectedScannedProps {
		actual := mock.rows.scanner.props[id]
		assert.IsType(expected, actual)
	}
}
