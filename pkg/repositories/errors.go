package repositories

import "github.com/Knoblauchpilze/backend-toolkit/pkg/errors"

const (
	errOptimisticLockException errors.ErrorCode = 200
)

var (
	ErrOptimisticLockException = errors.FromCode(errOptimisticLockException)
)
