package mappers

import "github.com/Knoblauchpilze/backend-toolkit/pkg/errors"

const (
	invalidEnumMapping errors.ErrorCode = 500
)

var (
	ErrInvalidEnumMapping = errors.FromCode(invalidEnumMapping)
)
