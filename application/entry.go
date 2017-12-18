package application

import (
	"strings"

	"github.com/pkg/errors"
	"github.com/russross/blackfriday"
	"github.com/takashabe/lumber/config"
	"github.com/takashabe/lumber/domain"
	"github.com/takashabe/lumber/domain/repository"
)

// EntryInteractor provides operation for entries
type EntryInteractor struct {
	repository repository.EntryRepository
}

// NewEntryInteractor returns initialized Entry object
func NewEntryInteractor(repo repository.EntryRepository) *EntryInteractor {
	return &EntryInteractor{repository: repo}
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

// Get returns entry when matched id
func (i *EntryInteractor) Get(id int) (*domain.Entry, error) {
	return i.repository.Get(id)
}

// Post saves the posted data in the background datastore
func (i *EntryInteractor) Post(e *EntryElement) (int, error) {
	if !e.Status.IsValid() {
		return 0, errors.Errorf("invalid entry status type: %d", e.Status)
	}
	return i.repository.Save(e.Entity())
}

// Edit changes entry the title and content
func (i *EntryInteractor) Edit(id int, e *EntryElement) error {
	entity := e.Entity()
	entity.ID = id
	return i.repository.Edit(entity)
}

// Delete deletes entry
func (i *EntryInteractor) Delete(id int) error {
	_, err := i.repository.Delete(id)
	return err
}

// EntryElement represent element of the entry operation method
type EntryElement struct {
	Title   string
	Content string
	Status  domain.EntryStatus
}

// NewEntryElement returns initialized an EntryElement object
func NewEntryElement(data []byte) (*EntryElement, error) {
	title, content := extractTitleAndContent(data)
	if len(title) == 0 || len(content) == 0 {
		return nil, config.ErrEmptyEntry
	}
	return &EntryElement{
		Title:   title,
		Content: content,
		Status:  domain.EntryStatusPublic,
	}, nil
}

// SetPublic set public status
func (e *EntryElement) SetPublic() {
	e.Status = domain.EntryStatusPublic
}

// SetPrivate set private status
func (e *EntryElement) SetPrivate() {
	e.Status = domain.EntryStatusPrivate
}

// Entity returns the entity from creating by the EntryElement
func (e *EntryElement) Entity() *domain.Entry {
	return &domain.Entry{
		Title:   e.Title,
		Content: e.Content,
		Status:  e.Status,
	}
}
