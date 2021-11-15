package errdef

import (
	"errors"
)

// Common errors used in ORAS
var (
	ErrNotFound      = errors.New("not found")
	ErrAlreadyExists = errors.New("already exists")
)
