package domainservices

import (
	"sort"
	"time"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
)

func AdvancePlanetToTime(
	planet *models.Planet,
	moment time.Time,
) error {
	events := buildEventsTimelineUntil(planet, moment)

	for _, event := range events {
		planet.UpdateToTime(event.completionTime)

		err := event.apply(planet)
		if err != nil {
			return err
		}
	}

	return planet.UpdateToTime(moment)
}

type completionEvent struct {
	completionTime time.Time
	apply          func(*models.Planet) error
}

func buildEventsTimelineUntil(planet *models.Planet, target time.Time) []completionEvent {
	var out []completionEvent

	if planet.BuildingAction != nil && planet.BuildingAction.CompletedAt.Compare(target) <= 0 {
		event := completionEvent{
			completionTime: planet.BuildingAction.CompletedAt,
			apply:          func(p *models.Planet) error { return p.ApplyBuildingAction() },
		}

		out = append(out, event)
	}

	for _, action := range planet.ShipActions {
		out = append(out, processShipAction(action, target)...)
	}

	sort.Slice(out, func(i, j int) bool {
		return out[i].completionTime.Before(out[j].completionTime)
	})

	return out
}

func processShipAction(action models.ShipAction, target time.Time) []completionEvent {
	var out []completionEvent

	unitCompleted := action.NextCompletionAt.Compare(target) <= 0

	for unitCompleted && !action.Completed() {
		event := completionEvent{
			completionTime: action.NextCompletionAt,
			apply:          func(p *models.Planet) error { return p.ApplyShipAction() },
		}
		action.CompleteOne()

		out = append(out, event)

		unitCompleted = action.NextCompletionAt.Compare(target) <= 0
	}

	return out
}
