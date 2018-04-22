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
	if _, ok := e.TrimPrivateTitle(); ok {
		e.Status = EntryStatusPrivate
	}
}

// TrimPrivateTitle returns contain private keyword in the title, and trimmed title
func (e *Entry) TrimPrivateTitle() (string, bool) {
	if len(e.Title) < 6 {
		return e.Title, false
	}

	var wip = []string{
		"[wip]",
		"[WIP]",
	}
	for _, w := range wip {
		if strings.HasPrefix(e.Title, w) {
			s := e.Title[len(w):]
			return strings.TrimSpace(s), true
		}
	}
	return e.Title, false
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
