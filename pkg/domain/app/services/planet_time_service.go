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

	// TODO: Maybe it would be nicer to have a slice of completionEvent
	// This slice is populated once based by checking the building action
	// and potentially other things (e.g. technology)
	// Each completion event is attached an apply function. This apply
	// function can either be a buildingApplier, or something else.
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
