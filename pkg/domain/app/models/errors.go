package models

import "github.com/Knoblauchpilze/backend-toolkit/pkg/errors"

var (
	noSuchBuilding errors.ErrorCode = 270

	errBuildingNotFound = errors.FromCode(noSuchBuilding)
)
