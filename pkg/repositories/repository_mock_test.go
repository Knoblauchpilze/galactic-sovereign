package repositories

import (
	"context"
	"time"

	"github.com/KnoblauchPilze/galactic-sovereign/pkg/db"
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

	timeStamp time.Time

	sqlQueries []string
	args       [][]interface{}

	queryCalled int
	execCalled  int

	rows         mockRows
	affectedRows []int
	execErrs     []error
}

func (m *mockTransaction) Close(ctx context.Context) {
	// We don't count the calls here because the repositories are not
	// closing the transaction. It is usually done at the service level.
}

func (m *mockTransaction) TimeStamp() time.Time {
	return m.timeStamp
}

func (m *mockTransaction) Query(ctx context.Context, sql string, arguments ...interface{}) db.Rows {
	m.queryCalled++

	if len(arguments) > 0 {
		var newArgs []interface{}
		newArgs = append(newArgs, arguments...)
		m.args = append(m.args, newArgs)
	}

	m.sqlQueries = append(m.sqlQueries, sql)

	return &m.rows
}

func (m *mockTransaction) Exec(ctx context.Context, sql string, arguments ...interface{}) (int, error) {
	var newArgs []interface{}
	newArgs = append(newArgs, arguments...)
	m.args = append(m.args, newArgs)

	m.sqlQueries = append(m.sqlQueries, sql)

	err := getValueToReturn(m.execCalled, m.execErrs)
	affectedRows := getValueToReturnOr(m.execCalled, m.affectedRows, 0)
	m.execCalled++

	if err == nil {
		return *affectedRows, nil
	}
	return *affectedRows, *err
}

func getValueToReturnOr[T any](count int, values []T, value T) *T {
	out := getValueToReturn(count, values)
	if out == nil {
		return &value
	}

	return out
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
