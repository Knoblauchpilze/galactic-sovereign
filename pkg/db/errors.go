package db

import (
	"github.com/KnoblauchPilze/user-service/pkg/errors"
)

const (
	NoMatchingSqlRows          errors.ErrorCode = 100
	MoreThanOneMatchingSqlRows errors.ErrorCode = 101
	OptimisticLockException    errors.ErrorCode = 102
	DuplicatedKeySqlKey        errors.ErrorCode = 103
	DatabasePingFailed         errors.ErrorCode = 110
)
