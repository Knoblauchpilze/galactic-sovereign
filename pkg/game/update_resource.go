package game

import (
	"time"

	"github.com/KnoblauchPilze/user-service/pkg/persistence"
)

func UpdateAmountToTime(resource persistence.PlanetResource, moment time.Time) persistence.PlanetResource {
	elapsed := moment.Sub(resource.UpdatedAt)
	if elapsed < 0 {
		return resource
	}

	hours := elapsed.Hours()
	resource.Amount += hours * float64(resource.Production)
	resource.UpdatedAt = moment

	return resource
}
