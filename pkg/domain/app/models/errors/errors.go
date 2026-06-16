package domainerrors

import "github.com/Knoblauchpilze/backend-toolkit/pkg/errors"

const (
	resourceNotFound        errors.ErrorCode = 600
	nameAlreadyTaken        errors.ErrorCode = 601
	actionAlreadyInProgress errors.ErrorCode = 602
)

var (
	ErrNotFound                = errors.FromCode(resourceNotFound)
	ErrNameAlreadyTaken        = errors.FromCode(nameAlreadyTaken)
	ErrActionAlreadyInProgress = errors.FromCode(actionAlreadyInProgress)
)
