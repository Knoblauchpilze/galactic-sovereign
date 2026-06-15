package domainerrors

import "github.com/Knoblauchpilze/backend-toolkit/pkg/errors"

const (
	resourceNotFound errors.ErrorCode = 600
)

var (
	ErrNotFound = errors.FromCode(resourceNotFound)
)
