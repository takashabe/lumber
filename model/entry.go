package model

import (
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/russross/blackfriday"
	"github.com/takashabe/lumber/config"
	"github.com/takashabe/lumber/datastore"
)

// Entry provides operation for entries
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

func (es EntryStatus) isValid() bool {
	return es.String() != "unknown"
}

// NewEntry returns initialized Entry object
func NewEntry(data []byte, status EntryStatus) (*Entry, error) {
	if !status.isValid() {
		return nil, errors.Errorf("invalid entry status type: %d", status)
	}

	title, content := extractTitleAndContent(data)
	if len(title) == 0 || len(content) == 0 {
		return nil, config.ErrEmptyEntry
	}
	return &Entry{
		Title:   title,
		Content: content,
		Status:  status,
	}, nil
}

func extractTitleAndContent(data []byte) (title, content string) {
	if len(data) == 0 {
		return "", ""
	}

	// Title is regarded as a first line
	titleIdx := 0
	html := blackfriday.Run(data)
	for i, d := range html {
		if d == '\n' {
			titleIdx = i
			break
		}
	}
	title = string(html[:titleIdx])
	title = trimHTMLTag(title)

	content = string(html[titleIdx:])
	content = strings.TrimSpace(content)
	return title, content
}

func trimHTMLTag(s string) string {
	openIdx := strings.IndexByte(s, '>')
	closeIdx := strings.LastIndexByte(s, '<')
	return s[openIdx+1 : closeIdx]
}

// GetEntry returns entry when matched id
func GetEntry(id int) (*Entry, error) {
	db, err := datastore.NewDatastore()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	raw, err := db.FindEntryByID(id)
	if err != nil {
		return nil, err
	}
	// TODO: consideration struct fields
	return &Entry{
		ID:        raw.ID,
		Title:     raw.Title,
		Content:   raw.Content,
		Status:    EntryStatus(raw.Status),
		CreatedAt: raw.CreatedAt,
		UpdatedAt: raw.UpdatedAt,
	}, nil
}

// Post saves the posted data in the background datastore
func (e *Entry) Post() error {
	db, err := datastore.NewDatastore()
	if err != nil {
		return err
	}
	defer db.Close()

	return db.SaveEntry(e.Title, e.Content, int(e.Status))
}
