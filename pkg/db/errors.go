package db

import (
	"github.com/KnoblauchPilze/galactic-sovereign/pkg/errors"
)

const (
	NoMatchingSqlRows          errors.ErrorCode = 100
	MoreThanOneMatchingSqlRows errors.ErrorCode = 101
	OptimisticLockException    errors.ErrorCode = 102
	DuplicatedKeySqlKey        errors.ErrorCode = 103
	DatabasePingFailed         errors.ErrorCode = 110
)
