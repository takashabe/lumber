package application

import (
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
	"github.com/takashabe/lumber/domain"
	"github.com/takashabe/lumber/domain/repository"
)

// TokenInteractor provides operation for tokens
type TokenInteractor struct {
	repository repository.TokenRepository
}

// NewTokenInteractor returns initialized token object
func NewTokenInteractor(repo repository.TokenRepository) *TokenInteractor {
	return &TokenInteractor{repository: repo}
}

// Get returns object when matched id
func (i *TokenInteractor) Get(id int) (*domain.Token, error) {
	return i.repository.Get(id)
}

// FindByValue returns object when matched value
func (i *TokenInteractor) FindByValue(v string) (*domain.Token, error) {
	return i.repository.FindByValue(v)
}

// New returns a token with a new unique value
func (i *TokenInteractor) New() (*domain.Token, error) {
	var maxAttempt = 20
	for a := 0; a < maxAttempt; a++ {
		v := generateToken()
		if _, err := i.FindByValue(v); err != nil {
			if err == domain.ErrNotFoundToken {
				token := &domain.Token{Value: v}
				return i.save(token)
			}
			return nil, err
		}
	}
	return nil, errors.New("failed to attempt for create a new token")
}

func (i *TokenInteractor) save(token *domain.Token) (*domain.Token, error) {
	id, err := i.repository.Save(token)
	if err != nil {
		return nil, err
	}
	token.ID = id
	return token, err
}

func generateToken() string {
	return uuid.NewV4().String()
}
