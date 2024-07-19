package repositories

import (
	"context"
	"fmt"
)

type SqlQueryType int

const (
	QueryBased SqlQueryType = 0
	ExecBased  SqlQueryType = 1
)

type testFunc[T any] func(context.Context, T) error
type testReturnFunc[T any] func(context.Context, T) interface{}

type dbInteractionTestCase[T any] struct {
	sqlMode SqlQueryType

	handler testFunc[T]

	expectedSql       string
	expectedArguments []interface{}
}

type dbReturnTestCase[T any] struct {
	handler         testReturnFunc[T]
	expectedContent interface{}
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

func getPoolCalledCount(sqlMode SqlQueryType, m *mockConnectionPoolNew) int {
	switch sqlMode {
	case QueryBased:
		return m.queryCalled
	case ExecBased:
		return m.execCalled
	default:
		panic(fmt.Errorf("Unsupported sql mode %v", sqlMode))
	}
}

func generatePoolErrorMock(sqlMode SqlQueryType, err error) *mockConnectionPoolNew {
	switch sqlMode {
	case QueryBased:
		return &mockConnectionPoolNew{
			rows: mockRowsNew{
				err: err,
			},
		}
	case ExecBased:
		return &mockConnectionPoolNew{
			execErr: err,
		}
	default:
		panic(fmt.Errorf("Unsupported sql mode %v", sqlMode))
	}
}

func getTransactionCalledCount(sqlMode SqlQueryType, m *mockTransactionNew) int {
	switch sqlMode {
	case QueryBased:
		return m.queryCalled
	case ExecBased:
		return m.execCalled
	default:
		panic(fmt.Errorf("Unsupported sql mode %v", sqlMode))
	}
}

func generateTransactionErrorMock(sqlMode SqlQueryType, err error) *mockTransactionNew {
	switch sqlMode {
	case QueryBased:
		return &mockTransactionNew{
			rows: mockRowsNew{
				err: err,
			},
		}
	case ExecBased:
		return &mockTransactionNew{
			execErr: err,
		}
	default:
		panic(fmt.Errorf("Unsupported sql mode %v", sqlMode))
	}
}
