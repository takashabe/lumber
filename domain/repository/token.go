package repository

import "github.com/takashabe/lumber/domain"

// TokenRepository represent reopsitory of the token
type TokenRepository interface {
	Get(id int) (*domain.Token, error)
	FindByValue(value string) (*domain.Token, error)
	Save(*domain.Token) (int, error)
	Update(*domain.Token) error
	Delete(id int) (bool, error)
}
