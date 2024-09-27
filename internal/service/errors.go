package service

import (
	"github.com/KnoblauchPilze/user-service/pkg/errors"
)

const (
	noSuchResource            errors.ErrorCode = 240
	FailedToCreateAction      errors.ErrorCode = 241
	ConflictingStateForAction errors.ErrorCode = 242
	ActionAlreadyCompleted    errors.ErrorCode = 243

	InvalidCredentials errors.ErrorCode = 250

	UserNotAuthenticated  errors.ErrorCode = 251
	AuthenticationExpired errors.ErrorCode = 252
)
