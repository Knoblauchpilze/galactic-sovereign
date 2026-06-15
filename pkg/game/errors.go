package game

import "github.com/Knoblauchpilze/backend-toolkit/pkg/errors"

const (
	errNoSuchBuilding     errors.ErrorCode = 270
	errNotEnoughResources errors.ErrorCode = 271

	errNoSuchResource errors.ErrorCode = 280

	errActionSchedulingFailed     errors.ErrorCode = 350
	errPlanetResourceUpdateFailed errors.ErrorCode = 351
)

var (
	ErrNoSuchBuilding     = errors.FromCode(errNoSuchBuilding)
	ErrNotEnoughResources = errors.FromCode(errNotEnoughResources)

	ErrNoSuchResource = errors.FromCode(errNoSuchResource)

	// ErrActionSchedulingFailed     errors.ErrorCode = 350
	// ErrPlanetResourceUpdateFailed errors.ErrorCode = 351
)
