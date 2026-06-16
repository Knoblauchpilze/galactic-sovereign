package service

import (
	"github.com/Knoblauchpilze/backend-toolkit/pkg/errors"
)

const (
	errInvalidEmail    errors.ErrorCode = 230
	errInvalidPassword errors.ErrorCode = 231

	errNoSuchResource               errors.ErrorCode = 240
	errActionUsesUnknownResource    errors.ErrorCode = 241
	errConflictingStateForAction    errors.ErrorCode = 242
	errActionAlreadyCompleted       errors.ErrorCode = 243
	errActionUpdatesUnknownResource errors.ErrorCode = 244

	errInvalidCredentials    errors.ErrorCode = 250
	errUserNotAuthenticated  errors.ErrorCode = 251
	errAuthenticationExpired errors.ErrorCode = 252
)

var (
	ErrInvalidEmail    = errors.FromCode(errInvalidEmail)
	ErrInvalidPassword = errors.FromCode(errInvalidPassword)

	ErrActionUsesUnknownResource    = errors.FromCode(errActionUsesUnknownResource)
	ErrConflictingStateForAction    = errors.FromCode(errConflictingStateForAction)
	ErrActionAlreadyCompleted       = errors.FromCode(errActionAlreadyCompleted)
	ErrActionUpdatesUnknownResource = errors.FromCode(errActionUpdatesUnknownResource)

	ErrInvalidCredentials    = errors.FromCode(errInvalidCredentials)
	ErrUserNotAuthenticated  = errors.FromCode(errUserNotAuthenticated)
	ErrAuthenticationExpired = errors.FromCode(errAuthenticationExpired)
)
