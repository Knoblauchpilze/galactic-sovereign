package drivenadapters

import (
	"context"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/db"
)

type DatabaseChecker struct {
	conn db.Connection
}

func NewDatabaseChecker(conn db.Connection) *DatabaseChecker {
	return &DatabaseChecker{
		conn: conn,
	}
}

func (d *DatabaseChecker) Ping(ctx context.Context) error {
	return d.conn.Ping(ctx)
}
