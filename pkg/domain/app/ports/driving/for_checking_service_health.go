package drivingports

import (
	"context"
)

type ForCheckingServiceHealth interface {
	Healthy(ctx context.Context) bool
}
