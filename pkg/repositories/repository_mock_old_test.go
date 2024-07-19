package repositories

import (
	"context"

	"github.com/KnoblauchPilze/user-service/pkg/db"
)

type mockConnectionPool struct {
	db.ConnectionPool

	queryCalled int
	execCalled  int

	affectedRows int
	execErr      error

	sqlQuery string
	args     []interface{}

	rows mockRows
}

func (m *mockConnectionPool) Query(ctx context.Context, sql string, arguments ...interface{}) db.Rows {
	m.queryCalled++
	m.sqlQuery = sql
	m.args = append(m.args, arguments...)
	return &m.rows
}

func (m *mockConnectionPool) Exec(ctx context.Context, sql string, arguments ...interface{}) (int, error) {
	m.execCalled++
	m.sqlQuery = sql
	m.args = append(m.args, arguments...)
	return m.affectedRows, m.execErr
}

type mockRows struct {
	err            error
	singleValueErr error
	allErr         error

	singleValueCalled int
	allCalled         int
	scanner           *mockScannable
}

func (m *mockRows) Err() error { return m.err }

func (m *mockRows) Empty() bool { return false }

func (m *mockRows) Close() {}

func (m *mockRows) GetSingleValue(parser db.RowParser) error {
	m.singleValueCalled++
	if m.scanner != nil {
		return parser(m.scanner)
	}
	return m.singleValueErr
}

func (m *mockRows) GetAll(parser db.RowParser) error {
	m.allCalled++
	if m.scanner != nil {
		return parser(m.scanner)
	}
	return m.allErr
}

type mockTransaction struct {
	db.Transaction

	queryCalled int
	execCalled  int

	affectedRows int
	execErr      error

	sqlQueries []string
	args       [][]interface{}

	rows mockRows
}

func (m *mockTransaction) Close(ctx context.Context) {}

func (m *mockTransaction) Query(ctx context.Context, sql string, arguments ...interface{}) db.Rows {
	m.queryCalled++

	var newArgs []interface{}
	newArgs = append(newArgs, arguments...)
	m.args = append(m.args, newArgs)

	m.sqlQueries = append(m.sqlQueries, sql)

	return &m.rows
}

func (m *mockTransaction) Exec(ctx context.Context, sql string, arguments ...interface{}) (int, error) {
	m.execCalled++

	var newArgs []interface{}
	newArgs = append(newArgs, arguments...)
	m.args = append(m.args, newArgs)

	m.sqlQueries = append(m.sqlQueries, sql)

	return m.affectedRows, m.execErr
}
