package drivenadapters

import (
	"github.com/Knoblauchpilze/backend-toolkit/pkg/db"
	domainerrors "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/errors"
)

func parseDbError(err error) error {
	if err == nil {
		return nil
	}

	if err == db.ErrNoMatchingRows {
		return domainerrors.ErrNotFound
	}

	if dbErr, ok := db.AsDatabaseError(err); ok {
		return parseFullDbError(dbErr)
	}

	return err
}

func parseFullDbError(err *db.DatabaseError) error {
	switch err.Code {
	case db.ErrForeignKeyValidation:
		return parseForeignKeyConstraintViolation(err)
	case db.ErrUniqueConstraintViolation:
		return parseUniqueConstraintViolation(err)
	default:
		return err
	}
}

func parseForeignKeyConstraintViolation(err *db.DatabaseError) error {
	switch err.Constraint {
	case "player_universe_fkey":
		return domainerrors.ErrUniverseNotFound
	default:
		return err
	}
}

func parseUniqueConstraintViolation(err *db.DatabaseError) error {
	switch err.Constraint {
	case "universe_name_key":
		return domainerrors.ErrNameAlreadyTaken
	case "player_universe_name_key":
		return domainerrors.ErrNameAlreadyTaken
	case "building_action_planet_key":
		return domainerrors.ErrActionAlreadyInProgress
	case "planet_coordinate_universe_galaxy_solar_system_position_key":
		return domainerrors.ErrCoordinateAlreadyUsed
	default:
		return err
	}
}
