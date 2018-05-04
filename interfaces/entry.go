package interfaces

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"

	"github.com/takashabe/lumber/application"
	"github.com/takashabe/lumber/config"
	"github.com/takashabe/lumber/domain"
	"github.com/takashabe/lumber/domain/repository"
)

// EntryHandler provides handler for the entry
type EntryHandler struct {
	entry *application.EntryInteractor
	auth  *application.AuthInteractor
}

// NewEntryHandler returns initialized EntryHandler
func NewEntryHandler(e repository.EntryRepository, t repository.TokenRepository) *EntryHandler {
	return &EntryHandler{
		entry: application.NewEntryInteractor(e),
		auth:  application.NewAuthInteractor(t),
	}
}

// Get returns entry when matched id
func (h *EntryHandler) Get(w http.ResponseWriter, r *http.Request, id int) {
	entry, err := h.entry.Get(id)
	if err != nil {
		Error(w, http.StatusNotFound, err, "failed to get entry")
		return
	}
	JSON(w, http.StatusOK, entry)
}

// GetIDs returns entry id list
func (h *EntryHandler) GetIDs(w http.ResponseWriter, r *http.Request) {
	ids, err := h.entry.GetIDs()
	if err != nil {
		Error(w, http.StatusNotFound, err, "failed to get entry")
		return
	}

	type response struct {
		IDs []int `json:"ids"`
	}
	JSON(w, http.StatusOK, response{IDs: ids})
}

// GetTitles returns entries
func (h *EntryHandler) GetTitles(w http.ResponseWriter, r *http.Request, start, length int) {
	es, err := h.entry.GetTitles(start, length)
	if err != nil {
		Error(w, http.StatusNotFound, err, "failed to get entry")
		return
	}

	type entry struct {
		ID    int    `json:"id"`
		Title string `json:"title"`
	}
	type response struct {
		Data []entry `json:"data"`
	}
	res := []entry{}
	for _, e := range es {
		res = append(res, entry{ID: e.ID, Title: e.Title})
	}
	JSON(w, http.StatusOK, response{Data: res})
}

// Post create new entry
func (h *EntryHandler) Post(w http.ResponseWriter, r *http.Request) {
	if err := h.authenticate(r); err != nil {
		Error(w, http.StatusUnauthorized, nil, "failed to authorized")
		return
	}

	raw := struct {
		Data   []byte `json:"data"`
		Status int    `json:"status"`
	}{}
	err := json.NewDecoder(r.Body).Decode(&raw)
	if err != nil {
		Error(w, http.StatusNotFound, err, "failed to parsed request")
		return
	}

	element, err := application.NewEntryElement(raw.Data)
	if err != nil {
		Error(w, http.StatusNotFound, err, "failed to create new entry")
		return
	}
	element.Status = domain.EntryStatus(raw.Status)
	id, err := h.entry.Post(element)
	if err != nil {
		if errors.Cause(err) == config.ErrDuplicatedTitle {
			JSON(w, http.StatusNoContent, nil)
			return
		}
		Error(w, http.StatusNotFound, err, "failed to create new entry")
		return
	}

	response := struct {
		ID int `json:"id"`
	}{
		ID: id,
	}
	JSON(w, http.StatusOK, response)
}

// Edit change entry the title and content
func (h *EntryHandler) Edit(w http.ResponseWriter, r *http.Request, id int) {
	if err := h.authenticate(r); err != nil {
		Error(w, http.StatusUnauthorized, nil, "failed to authorized")
		return
	}

	entry, err := h.entry.Get(id)
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

	element, err := application.NewEntryElement(raw.Data)
	if err != nil {
		Error(w, http.StatusNotFound, err, "failed to parse entry data")
		return
	}
	element.Status = entry.Status
	err = h.entry.Edit(id, element)
	if err != nil {
		Error(w, http.StatusNotFound, err, "failed to edit entry")
		return
	}
	JSON(w, http.StatusOK, nil)
}

// Delete deletes entry
func (h *EntryHandler) Delete(w http.ResponseWriter, r *http.Request, id int) {
	if err := h.authenticate(r); err != nil {
		Error(w, http.StatusUnauthorized, nil, "failed to authorized")
		return
	}

	_, err := h.entry.Get(id)
	if err != nil {
		Error(w, http.StatusNotFound, err, fmt.Sprintf("not found entry. id:%d", id))
		return
	}

	err = h.entry.Delete(id)
	if err != nil {
		Error(w, http.StatusNotFound, err, "failed to delete entry")
		return
	}
	JSON(w, http.StatusOK, nil)
}

func (h *EntryHandler) authenticate(r *http.Request) error {
	token := r.URL.Query().Get("token")
	if len(token) == 0 {
		return errors.New("invalid parmaeter")
	}
	return h.auth.AuthenticateByToken(token)
}
