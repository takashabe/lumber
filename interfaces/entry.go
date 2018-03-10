package interfaces

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/takashabe/lumber/application"
	"github.com/takashabe/lumber/domain"
	"github.com/takashabe/lumber/domain/repository"
)

// EntryHandler provides handler for the entry
type EntryHandler struct {
	interactor *application.EntryInteractor
}

// NewEntryHandler returns initialized EntryHandler
func NewEntryHandler(e repository.EntryRepository, t repository.TokenRepository) *EntryHandler {
	return &EntryHandler{
		interactor: application.NewEntryInteractor(e, t),
	}
}

// Get returns entry when matched id
func (h *EntryHandler) Get(w http.ResponseWriter, r *http.Request, id int) {
	entry, err := h.interactor.Get(id)
	if err != nil {
		Error(w, http.StatusNotFound, err, "failed to get entry")
		return
	}
	JSON(w, http.StatusOK, entry)
}

// GetIDs returns entry id list
func (h *EntryHandler) GetIDs(w http.ResponseWriter, r *http.Request) {
	ids, err := h.interactor.GetIDs()
	if err != nil {
		Error(w, http.StatusNotFound, err, "failed to get entry")
		return
	}

	type response struct {
		IDs []int `json:"ids"`
	}
	JSON(w, http.StatusOK, response{IDs: ids})
}

// Post create new entry
func (h *EntryHandler) Post(w http.ResponseWriter, r *http.Request) {
	token := getToken(r)
	if token == "" {
		Error(w, http.StatusBadRequest, nil, "invalid request parameters")
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
	id, err := h.interactor.Post(element, token)
	if err != nil {
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
	token := getToken(r)
	if token == "" {
		Error(w, http.StatusBadRequest, nil, "invalid request parameters")
		return
	}

	entry, err := h.interactor.Get(id)
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
	err = h.interactor.Edit(id, element, token)
	if err != nil {
		Error(w, http.StatusNotFound, err, "failed to edit entry")
		return
	}
	JSON(w, http.StatusOK, nil)
}

// Delete deletes entry
func (h *EntryHandler) Delete(w http.ResponseWriter, r *http.Request, id int) {
	token := getToken(r)
	if token == "" {
		Error(w, http.StatusBadRequest, nil, "invalid request parameters")
		return
	}

	_, err := h.interactor.Get(id)
	if err != nil {
		Error(w, http.StatusNotFound, err, fmt.Sprintf("not found entry. id:%d", id))
		return
	}

	err = h.interactor.Delete(id, token)
	if err != nil {
		Error(w, http.StatusNotFound, err, "failed to delete entry")
		return
	}
	JSON(w, http.StatusOK, nil)
}

func getToken(r *http.Request) string {
	return r.URL.Query().Get("token")
}
