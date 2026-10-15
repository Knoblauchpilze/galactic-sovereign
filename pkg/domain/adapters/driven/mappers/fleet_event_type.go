package mappers

import (
	"database/sql"

	domainerrors "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/errors"
)

type DbFleetEventType string

const (
	ArrivalEvent DbFleetEventType = "ARRIVAL"
	ReturnEvent  DbFleetEventType = "RETURN"
)

var (
	_ sql.Scanner = (*DbFleetEventType)(nil)

	allFleetEventTypes = map[DbFleetEventType]struct{}{
		ArrivalEvent: {},
		ReturnEvent:  {},
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

func validFleetEventType(event DbFleetEventType) bool {
	_, ok := allFleetEventTypes[event]
	return ok
}
