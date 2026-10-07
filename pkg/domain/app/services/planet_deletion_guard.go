package domainservices

import (
	"time"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	domainerrors "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/errors"
)

func PlanetDeletionGuard(planet *models.Planet, moment time.Time) error {
	err := AdvancePlanetToTime(planet, moment)
	if err != nil {
		return err
	}

	if planet.Homeworld {
		return domainerrors.ErrHomeworldCannotBeDeleted
	}

	if planet.BuildingAction != nil {
		return domainerrors.ErrBuildingActionNotCompleted
	}

	if len(planet.ShipActions) > 0 {
		return domainerrors.ErrShipActionNotCompleted
	}

	return nil
}
