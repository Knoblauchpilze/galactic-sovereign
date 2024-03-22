package db

import (
	"sync"

	"github.com/KnoblauchPilze/user-service/pkg/logger"
	"github.com/jackc/pgx"
)

type Connection interface {
	Connect() error
	Close()

	Query(sql string, arguments ...interface{}) Rows
	Exec(sql string, arguments ...interface{}) (string, error)
}

type pgxDbConnection interface {
	Close()
	Query(sql string, arguments ...interface{}) (*pgx.Rows, error)
	Exec(sql string, arguments ...interface{}) (pgx.CommandTag, error)
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

func (c *connectionImpl) Query(sql string, args ...interface{}) Rows {
	rows, err := c.pool.Query(sql, args...)
	return newRows(rows, err)
}

func (c *connectionImpl) Exec(sql string, args ...interface{}) (string, error) {
	tag, err := c.pool.Exec(sql, args...)
	return string(tag), err
}
