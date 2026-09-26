package infrastructure

import (
	"context"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/db"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/adapters/driven/database"
)

type dbConnectionImpl struct {
	*db.Connection
}

func NewDbConnection(conn *db.Connection) database.Connection {
	return &dbConnectionImpl{Connection: conn}
}

// Shadows the embedded BeginTx so the signature matches DbConnection.
func (c *dbConnectionImpl) BeginTx(ctx context.Context) (database.Transaction, error) {
	tx, err := c.Connection.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	return tx, nil
}
