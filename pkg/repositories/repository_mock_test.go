package repositories

import (
	"context"

	"github.com/KnoblauchPilze/user-service/pkg/db"
)

type mockConnectionPoolNew struct {
	db.ConnectionPool

	sqlQuery string
	args     []interface{}

	queryCalled int
	execCalled  int

	rows         mockRowsNew
	affectedRows int
	execErr      error
}

func (m *mockConnectionPoolNew) Query(ctx context.Context, sql string, arguments ...interface{}) db.Rows {
	m.sqlQuery = sql
	m.args = append(m.args, arguments...)

	m.queryCalled++

	return &m.rows
}

func (m *mockConnectionPoolNew) Exec(ctx context.Context, sql string, arguments ...interface{}) (int, error) {
	m.sqlQuery = sql
	m.args = append(m.args, arguments...)

	m.execCalled++

	return m.affectedRows, m.execErr
}

type mockRowsNew struct {
	db.Rows

	closeCalled          int
	getSingleValueCalled int
	getAllCalled         int

	err                error
	getSingleValueErrs []error
	getAllErrs         []error

	scanner *mockScannable
}

func (m *mockRowsNew) Err() error { return m.err }

func (m *mockRowsNew) Close() {
	m.closeCalled++
}

func (m *mockRowsNew) GetSingleValue(parser db.RowParser) error {
	if m.scanner != nil {
		return parser(m.scanner)
	}

	err := getValueToReturn(m.getSingleValueCalled, m.getSingleValueErrs)
	m.getSingleValueCalled++

	if err == nil {
		return nil
	}
	return *err
}

func (m *mockRowsNew) GetAll(parser db.RowParser) error {
	if m.scanner != nil {
		return parser(m.scanner)
	}

	err := getValueToReturn(m.getAllCalled, m.getAllErrs)
	m.getAllCalled++

	if err == nil {
		return nil
	}
	return *err
}

func getValueToReturn[T any](count int, values []T) *T {
	var out *T
	if count > len(values) {
		count = 0
	}
	if count < len(values) {
		out = &values[count]
	}

	return out
}
