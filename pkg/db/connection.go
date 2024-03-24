package db

import (
	"context"
	"sync"

	"github.com/KnoblauchPilze/user-service/pkg/logger"
	"github.com/jackc/pgx"
)

type Connection interface {
	Connect() error
	Close()

	Query(ctx context.Context, sql string, arguments ...interface{}) Rows
	Exec(ctx context.Context, sql string, arguments ...interface{}) (string, error)
}

type pgxDbConnection interface {
	Close()
	QueryEx(ctx context.Context, sql string, options *pgx.QueryExOptions, arguments ...interface{}) (*pgx.Rows, error)
	ExecEx(ctx context.Context, sql string, options *pgx.QueryExOptions, arguments ...interface{}) (pgx.CommandTag, error)
}

var pgxConnectionFunc = pgx.NewConnPool

type connectionImpl struct {
	config Config

	lock sync.Mutex
	pool pgxDbConnection
}

func NewConnection(config Config) Connection {
	return &connectionImpl{
		config: config,
	}
}

func (c *connectionImpl) Connect() error {
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

func (c *connectionImpl) Close() {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.pool == nil {
		return
	}

	c.pool.Close()
}

func (c *connectionImpl) Query(ctx context.Context, sql string, args ...interface{}) Rows {
	logger.Debugf("Query: %s (%d)", sql, len(args))
	rows, err := c.pool.QueryEx(ctx, sql, nil, args...)
	return newRows(rows, err)
}

func (c *connectionImpl) Exec(ctx context.Context, sql string, args ...interface{}) (string, error) {
	logger.Debugf("Exec: %s (%d)", sql, len(args))
	tag, err := c.pool.ExecEx(ctx, sql, nil, args...)
	return string(tag), err
}
