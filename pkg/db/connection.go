package db

import (
	"sync"

	"github.com/jackc/pgx"
)

type Connection interface {
	Connect() error
	Close()
}

var pgxConnectionFunc = pgx.NewConnPool

type connectionImpl struct {
	config Config

	lock sync.Mutex
	pool *pgx.ConnPool
}

func New(config Config) Connection {
	return &connectionImpl{
		config: config,
	}
}

func (c *connectionImpl) Connect() error {
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
