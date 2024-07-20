package repositories

import (
	"context"

	"github.com/KnoblauchPilze/user-service/pkg/db"
)

type mockConnectionPool struct {
	db.ConnectionPool

	sqlQuery string
	args     []interface{}

	queryCalled int
	execCalled  int

	rows         mockRows
	affectedRows int
	execErr      error
}

func (m *mockConnectionPool) Query(ctx context.Context, sql string, arguments ...interface{}) db.Rows {
	m.sqlQuery = sql
	m.args = append(m.args, arguments...)

	m.queryCalled++

	return &m.rows
}

func (m *mockConnectionPool) Exec(ctx context.Context, sql string, arguments ...interface{}) (int, error) {
	m.sqlQuery = sql
	m.args = append(m.args, arguments...)

	m.execCalled++

	return m.affectedRows, m.execErr
}

type mockRows struct {
	db.Rows

	closeCalled          int
	getSingleValueCalled int
	getAllCalled         int

	errs               []error
	getSingleValueErrs []error
	getAllErrs         []error

	scanner *mockScannable
}

func (m *mockRows) Err() error {
	err := getValueToReturn(m.getAllCalled+m.getSingleValueCalled, m.errs)

	if err == nil {
		return nil
	}
	return *err
}

// TODO: Add tests verifying this.
func (m *mockRows) Close() {
	m.closeCalled++
}

func (m *mockRows) GetSingleValue(parser db.RowParser) error {
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

func (m *mockRows) GetAll(parser db.RowParser) error {
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

type mockScannable struct {
	props [][]interface{}

	scanCalled int

	errs []error
}

func (m *mockScannable) Scan(dest ...interface{}) error {
	var newProps []interface{}
	newProps = append(newProps, dest...)
	m.props = append(m.props, newProps)

	err := getValueToReturn(m.scanCalled, m.errs)
	m.scanCalled++

	if err == nil {
		return nil
	}
	return *err
}

type mockTransaction struct {
	db.Transaction

	sqlQueries []string
	args       [][]interface{}

	closeCalled int
	queryCalled int
	execCalled  int

	rows         mockRows
	affectedRows int
	execErrs     []error
}

// TODO: Add tests to verify this
func (m *mockTransaction) Close(ctx context.Context) {
	m.closeCalled++
}

func (m *mockTransaction) Query(ctx context.Context, sql string, arguments ...interface{}) db.Rows {
	m.queryCalled++

	var newArgs []interface{}
	newArgs = append(newArgs, arguments...)
	m.args = append(m.args, newArgs)

	m.sqlQueries = append(m.sqlQueries, sql)

	return &m.rows
}

func (m *mockTransaction) Exec(ctx context.Context, sql string, arguments ...interface{}) (int, error) {
	var newArgs []interface{}
	newArgs = append(newArgs, arguments...)
	m.args = append(m.args, newArgs)

	m.sqlQueries = append(m.sqlQueries, sql)

	err := getValueToReturn(m.execCalled, m.execErrs)
	m.execCalled++

	if err == nil {
		return m.affectedRows, nil
	}
	return m.affectedRows, *err
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
