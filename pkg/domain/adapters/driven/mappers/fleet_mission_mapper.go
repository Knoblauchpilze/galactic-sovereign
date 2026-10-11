package mappers

import (
	"database/sql"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	domainerrors "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/errors"
)

type DbFleetMission models.FleetMission

var (
	_ sql.Scanner = (*DbFleetMission)(nil)

	allFleetMissions = map[DbFleetMission]struct{}{
		DbFleetMission(models.MissionColonize):  {},
		DbFleetMission(models.MissionTransport): {},
	}
)

func (mission *DbFleetMission) Scan(value any) error {
	if mission == nil {
		return domainerrors.ErrUnsupportedDatabaseEnumValue
	}

	var raw string
	switch value := value.(type) {
	case string:
		raw = value
	case []byte:
		raw = string(value)
	default:
		return domainerrors.ErrUnsupportedDatabaseEnumValue
	}

	parsed := DbFleetMission(raw)
	if !validFleetMission(parsed) {
		return domainerrors.ErrUnsupportedDatabaseEnumValue
	}

	*mission = parsed
	return nil
}

func (mission DbFleetMission) ToDomain() models.FleetMission {
	return models.FleetMission(mission)
}

func validFleetMission(mission DbFleetMission) bool {
	_, ok := allFleetMissions[mission]
	return ok
}
