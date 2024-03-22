package db

import (
	"testing"

	"github.com/KnoblauchPilze/user-service/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestRows_Err_NoError(t *testing.T) {
	assert := assert.New(t)

	r := newRows(nil, nil)
	assert.Nil(r.Err())
}

func TestRows_Err_SomeError(t *testing.T) {
	assert := assert.New(t)

	r := newRows(nil, errDefault)
	assert.Equal(errDefault, r.Err())
}

func TestRows_Close_DoesNotPanicWhenRowsIsNil(t *testing.T) {
	assert := assert.New(t)

	r := newRows(nil, nil)
	assert.NotPanics(r.Close)
}

func TestRows_Close_ClosesRows(t *testing.T) {
	assert := assert.New(t)

	m := mockSqlRows{}
	r := newRows(&m, nil)
	r.Close()
	assert.Equal(1, m.closeCalls)
}

type mockSqlRows struct {
	row        int
	rowsCount  int
	scanError  error
	closeCalls int
}

func (m *mockSqlRows) Next() bool {
	out := m.row < m.rowsCount
	m.row++
	return out
}

func (m *mockSqlRows) Scan(dest ...interface{}) error {
	return m.scanError
}

func (m *mockSqlRows) Close() {
	m.closeCalls++
}

func TestRows_Empty_NoRows(t *testing.T) {
	assert := assert.New(t)

	r := newRows(nil, nil)
	assert.True(r.Empty())
}

func TestRows_Empty_EmptyRows(t *testing.T) {
	assert := assert.New(t)

	mr := &mockSqlRows{}
	r := newRows(mr, nil)
	assert.True(r.Empty())
}

func TestRows_Empty_SomeRows(t *testing.T) {
	assert := assert.New(t)

	mr := &mockSqlRows{
		rowsCount: 1,
	}
	r := newRows(mr, nil)
	assert.False(r.Empty())
}

type mockParser struct {
	scanCalled int
	scanErr    error
}

func (m *mockParser) ScanRow(row Scannable) error {
	m.scanCalled++
	return m.scanErr
}

func TestRows_GetSingleValue_WhenError_Fails(t *testing.T) {
	assert := assert.New(t)

	mp := mockParser{}

	r := newRows(nil, errDefault)
	err := r.GetSingleValue(mp.ScanRow)
	assert.Equal(errDefault, err)
	assert.Equal(0, mp.scanCalled)
}

func TestRows_GetSingleValue_WhenNilRows_Fails(t *testing.T) {
	assert := assert.New(t)

	mp := &mockParser{}

	r := newRows(nil, nil)
	err := r.GetSingleValue(mp.ScanRow)
	assert.True(errors.IsErrorWithCode(err, NoMatchingSqlRows))
	assert.Equal(0, mp.scanCalled)
}

func TestRows_GetAll_WhenNoRows_Fails(t *testing.T) {
	assert := assert.New(t)

	mr := &mockSqlRows{}
	mp := &mockParser{}

	r := newRows(mr, nil)
	err := r.GetSingleValue(mp.ScanRow)

	assert.True(errors.IsErrorWithCode(err, NoMatchingSqlRows))
	assert.Equal(0, mr.closeCalls)
	assert.Equal(0, mp.scanCalled)
}

func TestRows_GetSingleValue_WhenRows_Succeeds(t *testing.T) {
	assert := assert.New(t)

	mr := &mockSqlRows{
		rowsCount: 1,
	}
	mp := &mockParser{}

	r := newRows(mr, nil)
	err := r.GetSingleValue(mp.ScanRow)

	assert.Nil(err)
	assert.Equal(1, mp.scanCalled)
}

func TestRows_GetSingleValue_CallsClose(t *testing.T) {
	assert := assert.New(t)

	mr := &mockSqlRows{
		rowsCount: 1,
	}
	mp := &mockParser{}

	r := newRows(mr, nil)
	r.GetSingleValue(mp.ScanRow)

	assert.Equal(1, mr.closeCalls)
}

func TestRows_GetSingleValue_ParserError(t *testing.T) {
	assert := assert.New(t)

	mr := &mockSqlRows{
		rowsCount: 1,
	}
	mp := &mockParser{
		scanErr: errDefault,
	}

	r := newRows(mr, nil)
	err := r.GetSingleValue(mp.ScanRow)

	assert.Equal(errDefault, err)
	assert.Equal(1, mr.closeCalls)
	assert.Equal(1, mp.scanCalled)
}

func TestRows_GetSingleValue_WithMultipleValues(t *testing.T) {
	assert := assert.New(t)

	mr := &mockSqlRows{
		rowsCount: 2,
	}
	mp := &mockParser{}

	r := newRows(mr, nil)
	err := r.GetSingleValue(mp.ScanRow)

	assert.True(errors.IsErrorWithCode(err, MoreThanOneMatchingSqlRows))
	assert.Equal(1, mr.closeCalls)
	assert.Equal(1, mp.scanCalled)
}

func TestRows_GetAll_WhenError_Fails(t *testing.T) {
	assert := assert.New(t)

	mp := &mockParser{}

	r := newRows(nil, errDefault)
	err := r.GetAll(mp.ScanRow)
	assert.Equal(errDefault, err)
	assert.Equal(0, mp.scanCalled)
}

func TestRows_GetAll_WhenNilRows_Succeeds(t *testing.T) {
	assert := assert.New(t)

	mp := &mockParser{}

	r := newRows(nil, nil)
	err := r.GetAll(mp.ScanRow)

	assert.Nil(err)
	assert.Equal(0, mp.scanCalled)
}

func TestRows_GetAll_WhenNoRows_Succeeds(t *testing.T) {
	assert := assert.New(t)

	mr := &mockSqlRows{}
	mp := &mockParser{}

	r := newRows(mr, nil)
	err := r.GetAll(mp.ScanRow)

	assert.Nil(err)
	assert.Equal(1, mr.closeCalls)
	assert.Equal(0, mp.scanCalled)
}

func TestRows_GetAll_WhenRows_Succeeds(t *testing.T) {
	assert := assert.New(t)

	mr := &mockSqlRows{
		rowsCount: 2,
	}
	mp := &mockParser{}

	r := newRows(mr, nil)
	err := r.GetAll(mp.ScanRow)

	assert.Nil(err)
	assert.Equal(2, mp.scanCalled)
	assert.Equal(1, mr.closeCalls)
}

func TestRows_GetAll_CallsClose(t *testing.T) {
	assert := assert.New(t)

	mr := &mockSqlRows{
		rowsCount: 1,
	}
	mp := &mockParser{}

	r := newRows(mr, nil)
	r.GetAll(mp.ScanRow)

	assert.Equal(1, mr.closeCalls)
}

func TestRows_GetAll_ParserError(t *testing.T) {
	assert := assert.New(t)

	mr := &mockSqlRows{
		rowsCount: 2,
	}
	mp := &mockParser{
		scanErr: errDefault,
	}

	r := newRows(mr, nil)
	err := r.GetAll(mp.ScanRow)

	assert.Equal(errDefault, err)
	assert.Equal(1, mr.closeCalls)
	assert.Equal(1, mp.scanCalled)
}
