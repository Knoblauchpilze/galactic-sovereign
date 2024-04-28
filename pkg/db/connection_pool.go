package db

import (
	"context"
	"sync"

	"github.com/KnoblauchPilze/user-service/pkg/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

type pgxConnectionFunc func(context.Context, *pgxpool.Config) (*pgxpool.Pool, error)

type connectionPoolImpl struct {
	config Config

	lock     sync.Mutex
	connFunc pgxConnectionFunc
	pool     pgxConnectionPool
}

func NewConnectionPool(config Config) ConnectionPool {
	return newConnectionPool(config, pgxpool.NewWithConfig)
}

func newConnectionPool(config Config, connFunc pgxConnectionFunc) ConnectionPool {
	return &connectionPoolImpl{
		config:   config,
		connFunc: connFunc,
	}
}

func (c *connectionPoolImpl) Connect(ctx context.Context) error {
	logger.Infof("Connecting to %s at %s:%d with user %s", c.config.Name, c.config.Host, c.config.Port, c.config.User)
	pool, err := c.connFunc(ctx, c.config.toConnPoolConfig())
	if err != nil {
		return err
	}

	c.lock.Lock()
	defer c.lock.Unlock()
	c.pool = pool

	return nil
}

func (c *connectionPoolImpl) Close() {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.pool == nil {
		return
	}

	c.pool.Close()
	logger.Infof("Closed connection to %s at %s:%d with user %s", c.config.Name, c.config.Host, c.config.Port, c.config.User)
}

func (c *connectionPoolImpl) StartTransaction(ctx context.Context) (Transaction, error) {
	pgxTx, err := c.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}

	tx := transactionImpl{
		tx: pgxTx,
	}
	return &tx, nil
}

func (c *connectionPoolImpl) Query(ctx context.Context, sql string, arguments ...any) Rows {
	log := logger.GetRequestLogger(ctx)
	log.Debugf("Query: %s (%d)", sql, len(arguments))

	rows, err := c.pool.Query(ctx, sql, arguments...)
	return newRows(rows, err)
}

func (c *connectionPoolImpl) Exec(ctx context.Context, sql string, arguments ...any) (int, error) {
	log := logger.GetRequestLogger(ctx)
	log.Debugf("Exec: %s (%d)", sql, len(arguments))

	tag, err := c.pool.Exec(ctx, sql, arguments...)
	return int(tag.RowsAffected()), err
}
