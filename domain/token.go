package domain

import "time"

// Token represent the token entity
type Token struct {
	ID        int
	Value     string
	CreatedAt time.Time
	UpdatedAt time.Time
}
