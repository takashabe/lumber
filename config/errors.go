package config

import "github.com/pkg/errors"

// Error constants
var (
	ErrInsufficientPrivileges = errors.New("insufficient privileges")
	ErrEmptyEntry             = errors.New("posting entry is empty")
	ErrEntrySizeLimitExceeded = errors.New("posting entry size is limit exceeded")
)
