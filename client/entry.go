package client

import (
	"context"
	"time"
)

// Entry provide operations associated with the entry
type Entry struct {
	// TODO: Depends specific API accessor service
	id    int
	addr  string
	token string
}

// EntryContent represent fields of the already published entry
type EntryContent struct {
	ID        int
	Title     string
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Get returns existed EntryContent
func (e *Entry) Get(ctx context.Context) (*EntryContent, error) {
	return nil, nil
}

// Edit submit makrdown file as an entry
func (e *Entry) Edit(ctx context.Context, file string) error {
	return nil
}

// Delete submit makrdown file as an entry
func (e *Entry) Delete(ctx context.Context) error {
	return nil
}
