package repositories

import (
	"context"

	"github.com/KnoblauchPilze/user-service/pkg/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type testSingleValueTxFunc func(context.Context, db.Transaction) error

type RepositorySingleValueTransactionTestSuite struct {
	suite.Suite

	testFunc testSingleValueTxFunc

	expectedSingleValueCalls int
	expectedScanCalls        int
	expectedScannedProps     [][]interface{}
}

func (s *RepositorySingleValueTransactionTestSuite) TestCallsGetSingleValue() {
	assert := assert.New(s.T())

	mock := &mockTransaction{}

	s.testFunc(context.Background(), mock)

	assert.Equal(s.expectedSingleValueCalls, mock.rows.singleValueCalled)
}

func (s *RepositorySingleValueTransactionTestSuite) TestPropagatesSingleValueError() {
	assert := assert.New(s.T())

	mock := &mockTransaction{
		rows: mockRows{
			singleValueErr: errDefault,
		},
	}

	err := s.testFunc(context.Background(), mock)

	assert.Equal(errDefault, err)
}

func (s *RepositorySingleValueTransactionTestSuite) TestPropagatesScanError() {
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

func (s *RepositorySingleValueTransactionTestSuite) TestWhenSingleValueSucceedsExpectsNoError() {
	assert := assert.New(s.T())

	mock := &mockTransaction{}

	err := s.testFunc(context.Background(), mock)

	assert.Nil(err)
}

func (s *RepositorySingleValueTransactionTestSuite) TestScansExpectedProperties() {
	assert := assert.New(s.T())

	mock := &mockTransaction{
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
