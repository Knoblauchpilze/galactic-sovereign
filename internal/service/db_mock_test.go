package service

import (
	"context"
	"time"

	"github.com/KnoblauchPilze/user-service/pkg/db"
)

type mockConnectionPool struct {
	db.ConnectionPool

	timeStamp time.Time

	txs  []*mockTransaction
	errs []error
}

func (m *mockConnectionPool) StartTransaction(ctx context.Context) (db.Transaction, error) {
	mt := &mockTransaction{
		timeStamp: m.timeStamp,
	}
	m.txs = append(m.txs, mt)

	var err error
	maybeErrorId := len(m.txs) - 1
	if maybeErrorId < len(m.errs) {
		err = m.errs[maybeErrorId]
	}

	return m.txs[len(m.txs)-1], err
}

type mockTransaction struct {
	db.Transaction

	timeStamp time.Time

	closeCalled int
}

func (m *mockTransaction) Close(ctx context.Context) {
	m.closeCalled++
}

func (m *mockTransaction) TimeStamp() time.Time {
	return m.timeStamp
}
