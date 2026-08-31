package domainerrors

import "github.com/Knoblauchpilze/backend-toolkit/pkg/errors"

const (
	resourceNotFound errors.ErrorCode = 600

	buildingNotFound       errors.ErrorCode = 602
	shipNotFound           errors.ErrorCode = 603
	universeNotFound       errors.ErrorCode = 604
	playerNotFound         errors.ErrorCode = 605
	planetResourceNotFound errors.ErrorCode = 606

	nameAlreadyTaken           errors.ErrorCode = 610
	actionAlreadyInProgress    errors.ErrorCode = 611
	noActionInProgress         errors.ErrorCode = 612
	notEnoughResources         errors.ErrorCode = 613
	optimisticLockingException errors.ErrorCode = 614
	planetNotUpToDate          errors.ErrorCode = 615
	buildingActionNotCompleted errors.ErrorCode = 616
	mutationWithoutVersionBump errors.ErrorCode = 617
	planetDeletionFailed       errors.ErrorCode = 618
	resourceCreationFailed     errors.ErrorCode = 619
	homeworldCannotBeDeleted   errors.ErrorCode = 620
	universeIsNotEmpty         errors.ErrorCode = 621
	coordinateAlreadyUsed      errors.ErrorCode = 622
	allFieldsUsed              errors.ErrorCode = 623
	shipActionNotCompleted     errors.ErrorCode = 624
)

var (
	ErrNotFound         = errors.FromCode(resourceNotFound)
	ErrBuildingNotFound = errors.FromCode(buildingNotFound)
	ErrShipNotFound     = errors.FromCode(shipNotFound)
	ErrUniverseNotFound = errors.FromCode(universeNotFound)
	ErrPlayerNotFound   = errors.FromCode(playerNotFound)
	ErrResourceNotFound = errors.FromCode(planetResourceNotFound)

	ErrNameAlreadyTaken           = errors.FromCode(nameAlreadyTaken)
	ErrActionAlreadyInProgress    = errors.FromCode(actionAlreadyInProgress)
	ErrNoActionInProgress         = errors.FromCode(noActionInProgress)
	ErrNotEnoughResources         = errors.FromCode(notEnoughResources)
	ErrOptimisticLocking          = errors.FromCode(optimisticLockingException)
	ErrPlanetNotUpToDate          = errors.FromCode(planetNotUpToDate)
	ErrBuildingActionNotCompleted = errors.FromCode(buildingActionNotCompleted)
	ErrMutationWithoutVersionBump = errors.FromCode(mutationWithoutVersionBump)
	ErrPlanetDeletionFailed       = errors.FromCode(planetDeletionFailed)
	ErrResourceCreationFailed     = errors.FromCode(resourceCreationFailed)
	ErrHomeworldCannotBeDeleted   = errors.FromCode(homeworldCannotBeDeleted)
	ErrUniverseIsNotEmpty         = errors.FromCode(universeIsNotEmpty)
	ErrCoordinateAlreadyUsed      = errors.FromCode(coordinateAlreadyUsed)
	ErrAllFieldsUsed              = errors.FromCode(allFieldsUsed)
	ErrShipActionNotCompleted     = errors.FromCode(shipActionNotCompleted)
)
