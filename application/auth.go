package application

import (
	"github.com/pkg/errors"
	"github.com/takashabe/lumber/config"
	"github.com/takashabe/lumber/domain/repository"
)

// AuthInteractor provides authentication
type AuthInteractor struct {
	tokenRepo repository.TokenRepository
}

// NewAuthInteractor returns initialized Auth object
func NewAuthInteractor(t repository.TokenRepository) *AuthInteractor {
	return &AuthInteractor{
		tokenRepo: t,
	}
}

// AuthenticateByToken provides validate of a token.
// Returns non-nil error when failed to authenticate.
func (i *AuthInteractor) AuthenticateByToken(repository repository.TokenRepository, token string) error {
	// TODO: Now process of the authenticate, only compare to exist a token.
	//       Want to add management of the user and authenticate each by user.
	_, err := repository.FindByValue(token)
	if err != nil {
		return errors.Wrapf(config.ErrInsufficientPrivileges, "error: %#v", err)
	}
	return nil
}
