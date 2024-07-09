package service

import (
	"github.com/KnoblauchPilze/user-service/pkg/errors"
)

const (
	InvalidCredentials errors.ErrorCode = 250

	UserNotAuthenticated  errors.ErrorCode = 251
	AuthenticationExpired errors.ErrorCode = 252
)
