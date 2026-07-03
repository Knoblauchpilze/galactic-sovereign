package drivenadapters

import (
	"context"
	"time"
)

type TimeAdapter struct{}

func NewTimeAdapter() *TimeAdapter {
	return &TimeAdapter{}
}

func (a *TimeAdapter) Now(_ context.Context) time.Time {
	return time.Now()
}
