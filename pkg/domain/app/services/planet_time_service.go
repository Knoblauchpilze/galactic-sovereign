package domainservices

import (
	"context"
	"time"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
)

func AdvancePlanetToTime(
	ctx context.Context,
	planet *models.Planet,
	moment time.Time,
) error {
	if planet.BuildingAction == nil {
		return planet.UpdateToTime(moment)
	}

	if planet.BuildingAction.CompletedAt.After(moment) {
		return planet.UpdateToTime(moment)
	}

	err := planet.UpdateToTime(planet.BuildingAction.CompletedAt)
	if err != nil {
		return err
	}

	err = planet.ApplyAction()
	if err != nil {
		return err
	}

	err = planet.UpdateToTime(moment)
	if err != nil {
		return err
	}

	return nil
}
