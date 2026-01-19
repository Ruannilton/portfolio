package auth

import (
	"errors"
)

var ErrUserNotFound = errors.New("user not found")
var ErrEmailAlreadyInUse = errors.New("email already in use")
var ErrFailedToGeneratePasswordHash = errors.New("failed to generate password hash")
var ErrFailedToGenerateAccessTokenHash = errors.New("failed to generate access token hash")
var ErrOAuthFailed = errors.New("oauth authentication failed")
