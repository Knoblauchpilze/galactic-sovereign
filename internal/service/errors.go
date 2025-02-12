package service

import (
	"github.com/Knoblauchpilze/backend-toolkit/pkg/errors"
)

const (
	InvalidEmail    errors.ErrorCode = 230
	InvalidPassword errors.ErrorCode = 231

	noSuchResource            errors.ErrorCode = 240
	ActionUsesUnknownResource errors.ErrorCode = 241
	ConflictingStateForAction errors.ErrorCode = 242
	ActionAlreadyCompleted    errors.ErrorCode = 243

	InvalidCredentials    errors.ErrorCode = 250
	UserNotAuthenticated  errors.ErrorCode = 251
	AuthenticationExpired errors.ErrorCode = 252
)
