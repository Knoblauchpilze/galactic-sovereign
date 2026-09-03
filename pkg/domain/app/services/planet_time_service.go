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

	if planet.BuildingAction != nil && planet.BuildingAction.CompletedAt.Before(target) {
		event := completionEvent{
			completionTime: planet.BuildingAction.CompletedAt,
			apply:          func(p *models.Planet) error { return p.ApplyBuildingAction() },
		}

		out = append(out, event)
	}

	for _, action := range planet.ShipActions {
		if action.NextCompletionAt.Before(target) {
			event := completionEvent{
				completionTime: action.NextCompletionAt,
				apply:          func(p *models.Planet) error { return p.ApplyShipAction() },
			}

			out = append(out, event)
		}
	}

	sort.Slice(out, func(i, j int) bool {
		return out[i].completionTime.Before(out[j].completionTime)
	})

	return out
}
