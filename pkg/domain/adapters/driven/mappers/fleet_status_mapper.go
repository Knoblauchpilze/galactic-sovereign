package mappers

import (
	"database/sql"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	domainerrors "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/errors"
)

type DbFleetStatus models.FleetStatus

var (
	_ sql.Scanner = (*DbFleetStatus)(nil)

	allFleetStatuses = map[DbFleetStatus]struct{}{
		DbFleetStatus(models.StatusOutbound):  {},
		DbFleetStatus(models.StatusReturning): {},
	}
)

func (status *DbFleetStatus) Scan(value any) error {
	if status == nil {
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

	parsed := DbFleetStatus(raw)
	if !validFleetStatus(parsed) {
		return domainerrors.ErrUnsupportedDatabaseEnumValue
	}

	*status = parsed
	return nil
}

func (status DbFleetStatus) ToDomain() models.FleetStatus {
	return models.FleetStatus(status)
}

func validFleetStatus(status DbFleetStatus) bool {
	_, ok := allFleetStatuses[status]
	return ok
}
