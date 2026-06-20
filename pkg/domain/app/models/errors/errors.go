package domainerrors

import "github.com/Knoblauchpilze/backend-toolkit/pkg/errors"

const (
	resourceNotFound errors.ErrorCode = 600

	buildingNotFound errors.ErrorCode = 602

	nameAlreadyTaken           errors.ErrorCode = 610
	actionAlreadyInProgress    errors.ErrorCode = 611
	noActionInProgress         errors.ErrorCode = 612
	notEnoughResources         errors.ErrorCode = 613
	optimisticLockingException errors.ErrorCode = 614
)

var (
	ErrNotFound         = errors.FromCode(resourceNotFound)
	ErrBuildingNotFound = errors.FromCode(buildingNotFound)

	ErrNameAlreadyTaken        = errors.FromCode(nameAlreadyTaken)
	ErrActionAlreadyInProgress = errors.FromCode(actionAlreadyInProgress)
	ErrNoActionInProgress      = errors.FromCode(noActionInProgress)
	ErrNotEnoughResources      = errors.FromCode(notEnoughResources)
	ErrOptimisticLocking       = errors.FromCode(optimisticLockingException)
)
