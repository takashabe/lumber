package domain

import (
	"github.com/pkg/errors"
)

// Token represent the token entity
type Token struct {
	ID    int
	Value string
}

// Token errors
var (
	ErrTokenAlreadyExistSameValue = errors.New("failed to save token. A record with the same value already exists")
	ErrNotFoundToken              = errors.New("failed to not found token")
)
