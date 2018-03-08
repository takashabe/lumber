package repository

import "github.com/takashabe/lumber/domain"

// EntryRepository represent reopsitory of the entry
type EntryRepository interface {
	Get(id int) (*domain.Entry, error)
	GetIDs() ([]int, error)
	Save(*domain.Entry) (int, error)
	Edit(*domain.Entry) error
	Delete(id int) (bool, error)
}
