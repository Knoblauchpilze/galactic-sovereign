package drivenadapters

import (
	"context"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/adapters/driven/database"
)

type DatabaseChecker struct {
	conn database.Connection
}

func NewDatabaseChecker(conn database.Connection) *DatabaseChecker {
	return &DatabaseChecker{
		conn: conn,
	}
}

func (d *DatabaseChecker) Ping(ctx context.Context) error {
	return d.conn.Ping(ctx)
}
