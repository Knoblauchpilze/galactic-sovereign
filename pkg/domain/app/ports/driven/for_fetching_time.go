package drivenports

import (
	"context"
	"time"
)

type ForFetchingTime interface {
	Now(ctx context.Context) time.Time
}
