package service

import (
	"context"

	"github.com/KnoblauchPilze/user-service/pkg/db"
)

type mockConnectionPool struct {
	db.ConnectionPool

	txs []*mockTransaction
	err error
}

func (m *mockConnectionPool) StartTransaction(ctx context.Context) (db.Transaction, error) {
	m.txs = append(m.txs, &mockTransaction{})

	return m.txs[len(m.txs)-1], m.err
}

type mockTransaction struct {
	db.Transaction

	closeCalled int
}

func (m *mockTransaction) Close(ctx context.Context) {
	m.closeCalled++
}
