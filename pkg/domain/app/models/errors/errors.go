package domainerrors

import "github.com/Knoblauchpilze/backend-toolkit/pkg/errors"

const (
	resourceNotFound errors.ErrorCode = 600

	buildingNotFound errors.ErrorCode = 602

	nameAlreadyTaken           errors.ErrorCode = 610
	actionAlreadyInProgress    errors.ErrorCode = 611
	notEnoughResources         errors.ErrorCode = 612
	optimisticLockingException errors.ErrorCode = 613
)

var (
	ErrNotFound         = errors.FromCode(resourceNotFound)
	ErrBuildingNotFound = errors.FromCode(buildingNotFound)

	ErrNameAlreadyTaken        = errors.FromCode(nameAlreadyTaken)
	ErrActionAlreadyInProgress = errors.FromCode(actionAlreadyInProgress)
	ErrNotEnoughResources      = errors.FromCode(notEnoughResources)
	ErrOptimisticLocking       = errors.FromCode(optimisticLockingException)
)
