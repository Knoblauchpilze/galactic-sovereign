package service

import (
	"context"

	"github.com/KnoblauchPilze/user-service/pkg/db"
)

type mockConnectionPool struct {
	db.ConnectionPool

	txs  []*mockTransaction
	errs []error
}

func (m *mockConnectionPool) StartTransaction(ctx context.Context) (db.Transaction, error) {
	m.txs = append(m.txs, &mockTransaction{})

	var err error
	maybeErrorId := len(m.txs) - 1
	if maybeErrorId < len(m.errs) {
		err = m.errs[maybeErrorId]
	}

	return m.txs[len(m.txs)-1], err
}

type mockTransaction struct {
	db.Transaction

	closeCalled int
}

func (m *mockTransaction) Close(ctx context.Context) {
	m.closeCalled++
}
