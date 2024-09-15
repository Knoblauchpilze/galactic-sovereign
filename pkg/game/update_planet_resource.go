package game

import (
	"time"

	"github.com/KnoblauchPilze/user-service/pkg/persistence"
)

func UpdatePlanetResourceAmountToTime(resource persistence.PlanetResource, production float64, moment time.Time) persistence.PlanetResource {
	elapsed := moment.Sub(resource.UpdatedAt)
	if elapsed < 0 {
		return resource
	}

	hours := elapsed.Hours()
	resource.Amount += hours * production
	resource.UpdatedAt = moment

	return resource
}
