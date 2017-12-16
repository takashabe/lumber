package application

import (
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/russross/blackfriday"
	"github.com/takashabe/lumber/config"
	"github.com/takashabe/lumber/infrastructure/persistence"
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
func NewEntry(data []byte) (*Entry, error) {
	title, content := extractTitleAndContent(data)
	if len(title) == 0 || len(content) == 0 {
		return nil, config.ErrEmptyEntry
	}
	return &Entry{
		Title:   title,
		Content: content,
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
	db, err := persistence.NewDatastore()
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
func (e *Entry) Post(s EntryStatus) (int, error) {
	if !s.isValid() {
		return 0, errors.Errorf("invalid entry status type: %d", s)
	}

	db, err := persistence.NewDatastore()
	if err != nil {
		return 0, err
	}
	defer db.Close()

	return db.SaveEntry(e.Title, e.Content, int(s))
// EntryElement represent element of the entry operation method
type EntryElement struct {
	Title   string
	Content string
	Status  domain.EntryStatus
}

// Edit changes entry the title and content
func (e *Entry) Edit() error {
	db, err := persistence.NewDatastore()
	if err != nil {
		return err
// NewEntryElement returns initialized an EntryElement object
func NewEntryElement(data []byte) (*EntryElement, error) {
	title, content := extractTitleAndContent(data)
	if len(title) == 0 || len(content) == 0 {
		return nil, config.ErrEmptyEntry
	}
	defer db.Close()
	return &EntryElement{
		Title:   title,
		Content: content,
		Status:  domain.EntryStatusPublic,
	}, nil
}

	return db.EditEntry(e.ID, e.Title, e.Content)
// SetPublic set public status
func (e *EntryElement) SetPublic() {
	e.Status = domain.EntryStatusPublic
}

// Delete deletes entry
func (e *Entry) Delete() error {
	db, err := persistence.NewDatastore()
	if err != nil {
		return err
	}
	defer db.Close()
// SetPrivate set private status
func (e *EntryElement) SetPrivate() {
	e.Status = domain.EntryStatusPrivate
}

	_, err = db.DeleteEntry(e.ID)
	return err
// Entity returns the entity from creating by the EntryElement
func (e *EntryElement) Entity() *domain.Entry {
	return &domain.Entry{
		Title:   e.Title,
		Content: e.Content,
		Status:  e.Status,
	}
}
