package db

import (
	"context"
	"sync"

	"github.com/KnoblauchPilze/user-service/pkg/logger"
	"github.com/KnoblauchPilze/user-service/pkg/middleware"
	"github.com/jackc/pgx"
)

type ConnectionPool interface {
	Connect() error
	Close()

	BeginTransaction(ctx context.Context) (Transaction, error)
	Query(ctx context.Context, sql string, arguments ...interface{}) Rows
	Exec(ctx context.Context, sql string, arguments ...interface{}) (int, error)
}

type pgxDbConnectionPool interface {
	Close()
	AcquireEx(ctx context.Context) (*pgx.Conn, error)
	QueryEx(ctx context.Context, sql string, options *pgx.QueryExOptions, arguments ...interface{}) (*pgx.Rows, error)
	ExecEx(ctx context.Context, sql string, options *pgx.QueryExOptions, arguments ...interface{}) (pgx.CommandTag, error)
}

var pgxConnectionFunc = pgx.NewConnPool
var pgxTransactionFunc = newTransactionFromPool

type connectionPoolImpl struct {
	config Config

	lock sync.Mutex
	pool pgxDbConnectionPool
}

func NewConnectionPool(config Config) ConnectionPool {
	return &connectionPoolImpl{
		config: config,
	}
}

func (c *connectionPoolImpl) Connect() error {
	logger.Infof("Connecting to %s at %s:%d with user %s", c.config.Name, c.config.Host, c.config.Port, c.config.User)
	pool, err := pgxConnectionFunc(c.config.toConnPoolConfig())
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

func (c *connectionPoolImpl) BeginTransaction(ctx context.Context) (Transaction, error) {
	return pgxTransactionFunc(ctx, c.pool)
}

func (c *connectionPoolImpl) Query(ctx context.Context, sql string, arguments ...interface{}) Rows {
	log := middleware.GetLoggerFromContext(ctx)
	log.Debugf("Query: %s (%d)", sql, len(arguments))

	rows, err := c.pool.QueryEx(ctx, sql, nil, arguments...)
	return newRows(rows, err)
}

func (c *connectionPoolImpl) Exec(ctx context.Context, sql string, arguments ...interface{}) (int, error) {
	log := middleware.GetLoggerFromContext(ctx)
	log.Debugf("Exec: %s (%d)", sql, len(arguments))

	tag, err := c.pool.ExecEx(ctx, sql, nil, arguments...)
	return int(tag.RowsAffected()), err
}
