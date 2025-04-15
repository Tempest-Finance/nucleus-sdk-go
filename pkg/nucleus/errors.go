package nucleus

import (
	"errors"
)

var (
	ErrFailedToExecute     = errors.New("failed to execute")
	ErrStrategiesIsInvalid = errors.New("strategies is invalid")
	ErrInvalidSigner       = errors.New("invalid signer")
	ErrEmptyCalls          = errors.New("empty calls")
	ErrStrategistNotFound  = errors.New("strategist not found")
)
