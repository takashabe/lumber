package interfaces

import (
	"net/http"

	"github.com/takashabe/lumber/application"
	"github.com/takashabe/lumber/domain/repository"
)

// TokenHandler provides handler for the token
type TokenHandler struct {
	interactor *application.TokenInteractor
}

// NewTokenHandler returns initialized TokenHandler
func NewTokenHandler(repo repository.TokenRepository) *TokenHandler {
	return &TokenHandler{
		interactor: application.NewTokenInteractor(repo),
	}
}

// Get returns token when mached id
func (h *TokenHandler) Get(w http.ResponseWriter, r *http.Request, id int) {
	token, err := h.interactor.Get(id)
	if err != nil {
		Error(w, http.StatusNotFound, err, "failed to get token")
		return
	}
	JSON(w, http.StatusOK, token)
}

// FindByValue returns token when mached value
func (h *TokenHandler) FindByValue(w http.ResponseWriter, r *http.Request, value string) {
	token, err := h.interactor.FindByValue(value)
	if err != nil {
		Error(w, http.StatusNotFound, err, "failed to get token")
		return
	}
	JSON(w, http.StatusOK, token)
}

// New returns a generated token
func (h *TokenHandler) New(w http.ResponseWriter, r *http.Request) {
	token, err := h.interactor.New()
	if err != nil {
		Error(w, http.StatusNotFound, err, "failed to create token")
		return
	}
	JSON(w, http.StatusCreated, token)
}
