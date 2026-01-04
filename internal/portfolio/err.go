package portfolio

import "errors"

var ErrProfileNotFound = errors.New("profile not found")
var ErrProfileAlreadyExists = errors.New("profile already exists for this user")
var ErrInvalidProfileData = errors.New("invalid profile data")