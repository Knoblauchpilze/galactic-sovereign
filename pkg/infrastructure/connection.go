package infrastructure

import (
	"context"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/db"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/adapters/driven/database"
)

// dbConnectionImpl provides a thin wrapper between the implementation
// of the connection from the backend-toolkit and the interface consumed
// by driven adapters.
// It is meant to be thin: if it grows too big it probably makes sense
// to revisit the interfaces in the driven adapters.
type dbConnectionImpl struct {
	*db.Connection
}

func NewDbConnection(conn *db.Connection) database.Connection {
	return &dbConnectionImpl{Connection: conn}
}

// Shadows the embedded BeginTx so the signature matches database.Connection.
func (c *dbConnectionImpl) BeginTx(ctx context.Context) (database.Transaction, error) {
	tx, err := c.Connection.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	return tx, nil
}
