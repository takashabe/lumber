package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
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
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Get returns existed EntryContent
func (e *Entry) Get(ctx context.Context) (*EntryContent, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%sapi/entry/%d", e.addr, e.id), nil)
	if err != nil {
		return nil, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	err = verifyHTTPStatusCode(http.StatusOK, res)
	if err != nil {
		return nil, err
	}

	buf := &EntryContent{}
	err = json.NewDecoder(res.Body).Decode(buf)
	return buf, err
}

// Edit submit makrdown file as an entry
func (e *Entry) Edit(ctx context.Context, file string) error {
	return nil
}

// Delete submit makrdown file as an entry
func (e *Entry) Delete(ctx context.Context) error {
	return nil
}
