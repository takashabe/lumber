package server

import (
	"net/http"

	"github.com/takashabe/lumber/model"
)

// GetEntry returns entry when matched id
func (s *Server) GetEntry(w http.ResponseWriter, r *http.Request, id int) {
	entry, err := model.GetEntry(id)
	if err != nil {
		Error(w, http.StatusNotFound, err, "failed to get entry")
	}
	JSON(w, http.StatusOK, entry)
}
