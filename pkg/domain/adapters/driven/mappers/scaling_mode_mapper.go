package mappers

import (
	"database/sql"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	domainerrors "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/errors"
)

type DbScalingMode models.ScalingMode

var _ sql.Scanner = (*DbScalingMode)(nil)

func (mode *DbScalingMode) Scan(value any) error {
	if mode == nil {
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

	parsed := DbScalingMode(raw)
	if parsed != DbScalingMode(models.LinearScaling) && parsed != DbScalingMode(models.GeometricScaling) {
		return domainerrors.ErrUnsupportedDatabaseEnumValue
	}

	*mode = parsed
	return nil
}

func (mode DbScalingMode) ToDomain() models.ScalingMode {
	return models.ScalingMode(mode)
}
