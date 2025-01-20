package game

import "github.com/Knoblauchpilze/backend-toolkit/pkg/errors"

const (
	NoSuchBuilding     errors.ErrorCode = 270
	NotEnoughResources errors.ErrorCode = 271

	NoSuchResource errors.ErrorCode = 280

	actionSchedulingFailed     errors.ErrorCode = 350
	planetResourceUpdateFailed errors.ErrorCode = 351
)
