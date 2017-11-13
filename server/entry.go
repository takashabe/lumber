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
		data   []byte
		status int
	}{}
	err := json.NewDecoder(r.Body).Decode(&raw)
	if err != nil {
		Error(w, http.StatusNotFound, err, "failed to parsed request")
		return
	}

	entry, err := model.NewEntry(raw.data, model.EntryStatus(raw.status))
	if err != nil {
		Error(w, http.StatusNotFound, err, "failed to create new entry")
		return
	}
	response := struct {
		id int
	}{
		id: entry.ID,
	}
	JSON(w, http.StatusOK, response)
}
