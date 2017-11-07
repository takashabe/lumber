package model

import (
	"strings"

	"github.com/pkg/errors"
	"github.com/russross/blackfriday"
	"github.com/takashabe/lumber/config"
	"github.com/takashabe/lumber/datastore"
)

// Entry provides operation for entries
type Entry struct {
	title   string
	content string
	status  EntryStatus
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
		title:   title,
		content: content,
		status:  status,
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

// Post saves the posted data in the background datastore
func (e *Entry) Post() error {
	db, err := datastore.NewDatastore()
	if err != nil {
		return err
	}
	defer db.Close()

	return db.SaveEntry(e.title, e.content, int(e.status))
}
