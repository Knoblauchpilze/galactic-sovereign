package mappers

import (
	"database/sql"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	domainerrors "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/errors"
)

type DbFleetEventType models.FleetEventType

var (
	_ sql.Scanner = (*DbFleetEventType)(nil)

	allFleetEventTypes = map[DbFleetEventType]struct{}{
		DbFleetEventType(models.ArrivalEvent): {},
		DbFleetEventType(models.ReturnEvent):  {},
	}
)

func (event *DbFleetEventType) Scan(value any) error {
	if event == nil {
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

	parsed := DbFleetEventType(raw)
	if !validFleetEventType(parsed) {
		return domainerrors.ErrUnsupportedDatabaseEnumValue
	}

	*event = parsed
	return nil
}

func (event DbFleetEventType) ToDomain() models.FleetEventType {
	return models.FleetEventType(event)
}

func validFleetEventType(event DbFleetEventType) bool {
	_, ok := allFleetEventTypes[event]
	return ok
}
