package game

import (
	"context"
	"time"

	"github.com/KnoblauchPilze/user-service/pkg/db"
	"github.com/google/uuid"
)

type PlanetResourceService interface {
	UpdatePlanetUntil(ctx context.Context, tx db.Transaction, planet uuid.UUID, until time.Time) error
}
