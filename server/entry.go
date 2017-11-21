package server

import (
	"encoding/json"
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

	entry, err := model.NewEntry(raw.Data, model.EntryStatus(raw.Status))
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
