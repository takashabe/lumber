package domain

import "strings"

// Entry represent the entry entity
type Entry struct {
	ID      int         `json:"id"`
	Title   string      `json:"title"`
	Content string      `json:"content"`
	Status  EntryStatus `json:"status"`
}

// UpdateStatusByTitle update entry status by title
// if has "[wip]" prefix, entry status to private
func (e *Entry) UpdateStatusByTitle() {
	if len(e.Title) < 6 {
		return
	}

	var wip = []string{
		"[wip]",
		"[WIP]",
	}
	for _, w := range wip {
		if strings.HasPrefix(e.Title, w) {
			e.Status = EntryStatusPrivate
			return
		}
	}
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

// IsValid returns whether the EntryStatus is valid
func (es EntryStatus) IsValid() bool {
	return es.String() != "unknown"
}
