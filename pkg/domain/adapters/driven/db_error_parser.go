package drivenadapters

import (
	"github.com/Knoblauchpilze/backend-toolkit/pkg/db"
	"github.com/Knoblauchpilze/backend-toolkit/pkg/errors"
	domainerrors "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/errors"
)

func parseDbError(err error) error {
	if err == nil {
		return nil
	}

	if errors.IsErrorWithCode(err, db.NoMatchingRows) {
		return domainerrors.ErrNotFound
	}

	return err
}
