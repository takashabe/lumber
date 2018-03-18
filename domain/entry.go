package domain

// Entry represent the entry entity
type Entry struct {
	ID      int         `json:"id"`
	Title   string      `json:"title"`
	Content string      `json:"content"`
	Status  EntryStatus `json:"status"`
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
