package repositories

import (
	"context"
	"fmt"

	"github.com/stretchr/testify/require"
)

type SqlQueryType int

const (
	QueryBased SqlQueryType = 0
	ExecBased  SqlQueryType = 1
)

type testFunc[T any] func(context.Context, T) error
type testReturnFunc[T any] func(context.Context, T) interface{}
type generateMock[T any] func() T
type verifyError func(error, *require.Assertions)

type dbInteractionTestCase[T any] struct {
	sqlMode      SqlQueryType
	generateMock generateMock[T]

	handler testFunc[T]

	expectedSqlQueries []string
	expectedArguments  [][]interface{}
}

type dbSingleValueTestCase[T any] struct {
	handler testFunc[T]

	expectedGetSingleValueCalls int
	expectedScanCalls           int
	expectedScannedProps        [][]interface{}
}

type dbGetAllTestCase[T any] struct {
	handler testFunc[T]

	expectedGetAllCalls  int
	expectedScanCalls    int
	expectedScannedProps [][]interface{}
}

type dbReturnTestCase[T any] struct {
	handler         testReturnFunc[T]
	expectedContent interface{}
}

type dbErrorTestCase[T any] struct {
	generateMock  generateMock[T]
	handler       testFunc[T]
	expectedError error
	verifyError   verifyError
}

func getPoolCalledCount(sqlMode SqlQueryType, m *mockConnectionPool) int {
	switch sqlMode {
	case QueryBased:
		return m.queryCalled
	case ExecBased:
		return m.execCalled
	default:
		panic(fmt.Errorf("Unsupported sql mode %v", sqlMode))
	}
}

func generatePoolErrorMock(sqlMode SqlQueryType, err error) *mockConnectionPool {
	switch sqlMode {
	case QueryBased:
		return &mockConnectionPool{
			rows: mockRows{
				errs: []error{err},
			},
		}
	case ExecBased:
		return &mockConnectionPool{
			execErr: err,
		}
	default:
		panic(fmt.Errorf("Unsupported sql mode %v", sqlMode))
	}
}

func getTransactionCalledCount(sqlMode SqlQueryType, m *mockTransaction) int {
	switch sqlMode {
	case QueryBased:
		return m.queryCalled
	case ExecBased:
		return m.execCalled
	default:
		panic(fmt.Errorf("Unsupported sql mode %v", sqlMode))
	}
}

func generateTransactionErrorMock(sqlMode SqlQueryType, err error) *mockTransaction {
	switch sqlMode {
	case QueryBased:
		return &mockTransaction{
			rows: mockRows{
				errs: []error{err},
			},
		}
	case ExecBased:
		return &mockTransaction{
			execErrs: []error{err},
		}
	default:
		panic(fmt.Errorf("Unsupported sql mode %v", sqlMode))
	}
}
