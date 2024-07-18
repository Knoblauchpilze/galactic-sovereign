package service

import (
	"context"

	"github.com/KnoblauchPilze/user-service/pkg/db"
)

type mockConnectionPool struct {
	db.ConnectionPool

	tx  mockTransaction
	err error
}

func (m *mockConnectionPool) StartTransaction(ctx context.Context) (db.Transaction, error) {
	return &m.tx, m.err
}

type mockTransaction struct {
	db.Transaction

	closeCalled int
}

func (m *mockTransaction) Close(ctx context.Context) {
	m.closeCalled++
}
