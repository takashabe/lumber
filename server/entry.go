package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/takashabe/lumber/model"
)

// GetEntry returns entry when matched id
func (s *Server) GetEntry(w http.ResponseWriter, r *http.Request, id int) {
	entry, err := model.GetEntry(id)
	if err != nil {
		Error(w, http.StatusNotFound, err, "failed to get entry")
		return
	}
	JSON(w, http.StatusOK, entry)
}

// PostEntry create new entry
func (s *Server) PostEntry(w http.ResponseWriter, r *http.Request) {
	raw := struct {
		Data   []byte `json:"data"`
		Status int    `json:"status"`
	}{}
	err := json.NewDecoder(r.Body).Decode(&raw)
	if err != nil {
		Error(w, http.StatusNotFound, err, "failed to parsed request")
		return
	}

	entry, err := model.NewEntry(raw.Data)
	if err != nil {
		Error(w, http.StatusNotFound, err, "failed to create new entry")
		return
	}
	err = entry.Post(model.EntryStatus(raw.Status))
	if err != nil {
		Error(w, http.StatusNotFound, err, "failed to create new entry")
		return
	}
	response := struct {
		ID int `json:"id"`
	}{
		ID: entry.ID,
	}
	JSON(w, http.StatusOK, response)
}

// EditEntry change entry the title and content
func (s *Server) EditEntry(w http.ResponseWriter, r *http.Request, id int) {
	_, err := model.GetEntry(id)
	if err != nil {
		Error(w, http.StatusNotFound, err, fmt.Sprintf("not found entry. id:%d", id))
		return
	}

	raw := struct {
		Data []byte `json:"data"`
	}{}
	err = json.NewDecoder(r.Body).Decode(&raw)
	if err != nil {
		Error(w, http.StatusNotFound, err, "failed to parse request")
		return
	}

	entry, err := model.NewEntry(raw.Data)
	if err != nil {
		Error(w, http.StatusNotFound, err, "failed to parse entry data")
		return
	}
	entry.ID = id
	err = entry.Edit()
	if err != nil {
		Error(w, http.StatusNotFound, err, "failed to edit entry")
		return
	}
	JSON(w, http.StatusOK, nil)
}

// DeleteEntry deletes entry
func (s *Server) DeleteEntry(w http.ResponseWriter, r *http.Request, id int) {
	entry, err := model.GetEntry(id)
	if err != nil {
		Error(w, http.StatusNotFound, err, fmt.Sprintf("not found entry. id:%d", id))
		return
	}
	err = entry.Delete()
	if err != nil {
		Error(w, http.StatusNotFound, err, "failed to delete entry")
		return
	}
	JSON(w, http.StatusOK, nil)
}
