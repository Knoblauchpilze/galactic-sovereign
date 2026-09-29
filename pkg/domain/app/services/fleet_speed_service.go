package domainservices

import (
	"math"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	domainerrors "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/errors"
	"github.com/google/uuid"
)

func DetermineFleetSpeed(
	order models.FleetOrder,
	ships []models.Ship,
) (int, error) {
	if len(order.Ships) == 0 {
		return 0, domainerrors.ErrNoShipInFleet
	}

	temp := make(map[uuid.UUID]int)
	for _, s := range ships {
		temp[s.Id] = s.BaseSpeed
	}

	lowestSpeed := math.MaxInt

	for _, s := range order.Ships {
		shipSpeed, ok := temp[s.Ship]
		if !ok {
			return 0, domainerrors.ErrShipNotFound
		}

		if shipSpeed < lowestSpeed {
			lowestSpeed = shipSpeed
		}
	}

	return lowestSpeed, nil
}
