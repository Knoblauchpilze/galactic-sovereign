package db

import (
	"testing"

	"github.com/KnoblauchPilze/user-service/pkg/errors"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
)

func TestRows_Err_NoError(t *testing.T) {
	assert := assert.New(t)

	r := newRows(nil, nil)
	assert.Nil(r.Err())
}

func TestRows_Err_WithError(t *testing.T) {
	assert := assert.New(t)

	r := newRows(nil, errDefault)
	assert.Equal(errDefault, r.Err())
}

func TestRows_Close_DoesNotPanicWhenRowsIsNil(t *testing.T) {
	assert := assert.New(t)

	r := rowsImpl{}
	assert.NotPanics(r.close)
}

type mockPgxRows struct {
	pgx.Rows

	row         int
	rowsCount   int
	scanError   error
	closeCalled int
}

func TestRows_Close_ClosesRows(t *testing.T) {
	assert := assert.New(t)

	m := mockPgxRows{}
	r := rowsImpl{
		rows: &m,
	}
	r.close()
	assert.Equal(1, m.closeCalled)
}

type mockParser struct {
	scanCalled int
	err        error
}

func TestRows_GetSingleValue_WhenError_Fails(t *testing.T) {
	assert := assert.New(t)

	mp := &mockParser{}

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

	mr := &mockPgxRows{}
	mp := &mockParser{}

	r := newRows(mr, nil)
	err := r.GetSingleValue(mp.ScanRow)

	assert.True(errors.IsErrorWithCode(err, NoMatchingSqlRows))
	assert.Equal(0, mr.closeCalled)
	assert.Equal(0, mp.scanCalled)
}

func TestRows_GetSingleValue_CallsScan(t *testing.T) {
	assert := assert.New(t)

	mr := &mockPgxRows{
		rowsCount: 1,
	}
	mp := &mockParser{}

	r := newRows(mr, nil)
	r.GetSingleValue(mp.ScanRow)

	assert.Equal(1, mp.scanCalled)
}

func TestRows_GetSingleValue_ReturnsScanError(t *testing.T) {
	assert := assert.New(t)

	mr := &mockPgxRows{
		rowsCount: 1,
	}
	mp := &mockParser{
		err: errDefault,
	}

	r := newRows(mr, nil)
	err := r.GetSingleValue(mp.ScanRow)

	assert.Equal(errDefault, err)
}

func TestRows_GetSingleValue_CallsClose(t *testing.T) {
	assert := assert.New(t)

	mr := &mockPgxRows{
		rowsCount: 1,
	}
	mp := &mockParser{}

	r := newRows(mr, nil)
	r.GetSingleValue(mp.ScanRow)

	assert.Equal(1, mr.closeCalled)
}

func TestRows_GetSingleValue_CallsCloseAlsoWhenScanFails(t *testing.T) {
	assert := assert.New(t)

	mr := &mockPgxRows{
		rowsCount: 1,
	}
	mp := &mockParser{
		err: errDefault,
	}

	r := newRows(mr, nil)
	r.GetSingleValue(mp.ScanRow)

	assert.Equal(1, mp.scanCalled)
	assert.Equal(1, mr.closeCalled)
}

func TestRows_GetSingleValue_WithMultipleValues_Fails(t *testing.T) {
	assert := assert.New(t)

	mr := &mockPgxRows{
		rowsCount: 2,
	}
	mp := &mockParser{}

	r := newRows(mr, nil)
	err := r.GetSingleValue(mp.ScanRow)

	assert.True(errors.IsErrorWithCode(err, MoreThanOneMatchingSqlRows))
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

	mr := &mockPgxRows{}
	mp := &mockParser{}

	r := newRows(mr, nil)
	err := r.GetAll(mp.ScanRow)

	assert.Nil(err)
	assert.Equal(0, mp.scanCalled)
}

func TestRows_GetAll_WhenRows_Succeeds(t *testing.T) {
	assert := assert.New(t)

	mr := &mockPgxRows{
		rowsCount: 2,
	}
	mp := &mockParser{}

	r := newRows(mr, nil)
	err := r.GetAll(mp.ScanRow)

	assert.Nil(err)
	assert.Equal(2, mp.scanCalled)
}

func TestRows_GetAll_CallsClose(t *testing.T) {
	assert := assert.New(t)

	mr := &mockPgxRows{
		rowsCount: 1,
	}
	mp := &mockParser{}

	r := newRows(mr, nil)
	r.GetAll(mp.ScanRow)

	assert.Equal(1, mr.closeCalled)
}

func TestRows_GetAll_ReturnsScanError(t *testing.T) {
	assert := assert.New(t)

	mr := &mockPgxRows{
		rowsCount: 2,
	}
	mp := &mockParser{
		err: errDefault,
	}

	r := newRows(mr, nil)
	err := r.GetAll(mp.ScanRow)

	assert.Equal(errDefault, err)
	assert.Equal(1, mp.scanCalled)
}

func TestRows_GetAll_CallsCloseAlsoWhenScanFails(t *testing.T) {
	assert := assert.New(t)

	mr := &mockPgxRows{
		rowsCount: 1,
	}
	mp := &mockParser{
		err: errDefault,
	}

	r := newRows(mr, nil)
	r.GetAll(mp.ScanRow)

	assert.Equal(1, mp.scanCalled)
}

func (m *mockPgxRows) Next() bool {
	out := m.row < m.rowsCount
	m.row++
	return out
}

func (m *mockPgxRows) Scan(dest ...interface{}) error {
	return m.scanError
}

func (m *mockPgxRows) Close() {
	m.closeCalled++
}

func (m *mockParser) ScanRow(row Scannable) error {
	m.scanCalled++
	return m.err
}
