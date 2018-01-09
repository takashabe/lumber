package application

import "github.com/takashabe/lumber/domain/repository"

// authenticateByToken provides validate of a token.
// Returns non-nil error when failed to authenticate.
func authenticateByToken(repository *repository.TokenRepository, token string) error {
	// TODO: Now process of the authenticate, only compare to exist a token.
	//       Want to add management of the user and authenticate each by user.
	_, err := repository.FindByValue(token)
	return err
}
