package domain

import (
	"time"

	"github.com/pkg/errors"
)

// Token represent the token entity
type Token struct {
	ID        int
	Value     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Token errors
var (
	ErrTokenAlreadyExistSameValue = errors.New("failed to save token. A record with the same value already exists")
)
