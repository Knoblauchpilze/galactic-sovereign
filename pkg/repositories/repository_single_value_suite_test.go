package repositories

import (
	"context"

	"github.com/KnoblauchPilze/user-service/pkg/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type testSingleValueFunc func(context.Context, db.ConnectionPool) error

type RepositorySingleValueTestSuite struct {
	suite.Suite

	testFunc testSingleValueFunc

	expectedScanCalls    int
	expectedScannedProps []interface{}
}

func (s *RepositorySingleValueTestSuite) TestCallsGetSingleValue() {
	assert := assert.New(s.T())

	mock := &mockConnectionPool{}

	s.testFunc(context.Background(), mock)

	assert.Equal(1, mock.rows.singleValueCalled)
}

func (s *RepositorySingleValueTestSuite) TestPropagatesSingleValueError() {
	assert := assert.New(s.T())

	mock := &mockConnectionPool{
		rows: mockRows{
			singleValueErr: errDefault,
		},
	}

	err := s.testFunc(context.Background(), mock)

	assert.Equal(errDefault, err)
}

func (s *RepositorySingleValueTestSuite) TestPropagatesScanError() {
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

func (s *RepositorySingleValueTestSuite) TestWhenSingleValueSucceedsExpectsNoError() {
	assert := assert.New(s.T())

	mock := &mockConnectionPool{}

	err := s.testFunc(context.Background(), mock)

	assert.Nil(err)
}

func (s *RepositorySingleValueTestSuite) TestScansExpectedProperties() {
	assert := assert.New(s.T())

	mock := &mockConnectionPool{
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
