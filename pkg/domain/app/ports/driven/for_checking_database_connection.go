package drivenports

import (
	"context"
)

type ForCheckingDatabaseConnection interface {
	Ping(ctx context.Context) error
}
