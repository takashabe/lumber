package domain

import "time"

// Entry represent the entry entity
type Entry struct {
	ID        int
	Title     string
	Content   string
	Status    EntryStatus
	CreatedAt time.Time
	UpdatedAt time.Time
}

// EntryStatus represent status of the entries
type EntryStatus int

// EntryStatus details
const (
	EntryStatusPublic EntryStatus = iota
	EntryStatusPrivate
)

func (es EntryStatus) String() string {
	switch es {
	case EntryStatusPublic:
		return "public"
	case EntryStatusPrivate:
		return "private"
	default:
		return "unknown"
	}
}
