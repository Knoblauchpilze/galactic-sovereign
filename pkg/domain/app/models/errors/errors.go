package domainerrors

import "github.com/Knoblauchpilze/backend-toolkit/pkg/errors"

const (
	resourceNotFound errors.ErrorCode = 600

	buildingNotFound       errors.ErrorCode = 602
	universeNotFound       errors.ErrorCode = 603
	playerNotFound         errors.ErrorCode = 604
	planetResourceNotFound errors.ErrorCode = 605

	nameAlreadyTaken           errors.ErrorCode = 610
	actionAlreadyInProgress    errors.ErrorCode = 611
	noActionInProgress         errors.ErrorCode = 612
	notEnoughResources         errors.ErrorCode = 613
	optimisticLockingException errors.ErrorCode = 614
	planetNotUpToDate          errors.ErrorCode = 615
	actionNotCompleted         errors.ErrorCode = 616
	mutationWithoutVersionBump errors.ErrorCode = 617
	planetDeletionFailed       errors.ErrorCode = 618
)

var (
	ErrNotFound         = errors.FromCode(resourceNotFound)
	ErrBuildingNotFound = errors.FromCode(buildingNotFound)
	ErrUniverseNotFound = errors.FromCode(universeNotFound)
	ErrPlayerNotFound   = errors.FromCode(playerNotFound)
	ErrResourceNotFound = errors.FromCode(planetResourceNotFound)

	ErrNameAlreadyTaken           = errors.FromCode(nameAlreadyTaken)
	ErrActionAlreadyInProgress    = errors.FromCode(actionAlreadyInProgress)
	ErrNoActionInProgress         = errors.FromCode(noActionInProgress)
	ErrNotEnoughResources         = errors.FromCode(notEnoughResources)
	ErrOptimisticLocking          = errors.FromCode(optimisticLockingException)
	ErrPlanetNotUpToDate          = errors.FromCode(planetNotUpToDate)
	ErrActionNotCompleted         = errors.FromCode(actionNotCompleted)
	ErrMutationWithoutVersionBump = errors.FromCode(mutationWithoutVersionBump)
	ErrPlanetDeletionFailed       = errors.FromCode(planetDeletionFailed)
)
