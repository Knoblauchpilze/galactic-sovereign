package drivenadapters

import (
	"context"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/db"
	drivenports "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/ports/driven"
)

type databaseChecker struct {
	conn db.Connection
}

func NewDatabaseChecker(conn db.Connection) drivenports.ForCheckingDatabaseConnection {
	return &databaseChecker{
		conn: conn,
	}
}

func (d *databaseChecker) Ping(ctx context.Context) error {
	return d.conn.Ping(ctx)
}
