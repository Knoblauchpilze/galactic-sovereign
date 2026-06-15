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

	return err
}
