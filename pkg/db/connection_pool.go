package db

import (
	"context"
	"sync"

	"github.com/KnoblauchPilze/user-service/pkg/logger"
	"github.com/jackc/pgx"
)

type pgxConnectionFunc func(pgx.ConnPoolConfig) (*pgx.ConnPool, error)

type connectionPoolImpl struct {
	config Config

	lock     sync.Mutex
	connFunc pgxConnectionFunc
	pool     pgxConnectionPool
}

func NewConnectionPool(config Config) ConnectionPool {
	return newConnectionPool(config, pgx.NewConnPool)
}

func newConnectionPool(config Config, connFunc pgxConnectionFunc) ConnectionPool {
	return &connectionPoolImpl{
		config:   config,
		connFunc: connFunc,
	}
}

func (c *connectionPoolImpl) Connect() error {
	logger.Infof("Connecting to %s at %s:%d with user %s", c.config.Name, c.config.Host, c.config.Port, c.config.User)
	pool, err := c.connFunc(c.config.toConnPoolConfig())
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
	pgxTx, err := c.pool.BeginEx(ctx, nil)
	if err != nil {
		return nil, err
	}

	tx := transactionImpl{
		tx: pgxTx,
	}
	return &tx, nil
}

func (c *connectionPoolImpl) Query(ctx context.Context, sql string, arguments ...interface{}) Rows {
	log := logger.GetRequestLogger(ctx)
	log.Debugf("Query: %s (%d)", sql, len(arguments))

	rows, err := c.pool.QueryEx(ctx, sql, nil, arguments...)
	return newRows(rows, err)
}

func (c *connectionPoolImpl) Exec(ctx context.Context, sql string, arguments ...interface{}) (int, error) {
	log := logger.GetRequestLogger(ctx)
	log.Debugf("Exec: %s (%d)", sql, len(arguments))

	tag, err := c.pool.ExecEx(ctx, sql, nil, arguments...)
	return int(tag.RowsAffected()), err
}
