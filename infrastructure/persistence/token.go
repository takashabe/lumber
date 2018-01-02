package persistence

import (
	"database/sql"

	"github.com/takashabe/lumber/domain"
	"github.com/takashabe/lumber/domain/repository"
	"github.com/takashabe/lumber/infrastructure/utils"
)

// TokenRepositoryImpl implements the TokenRepository
type TokenRepositoryImpl struct {
	*SQLRepositoryAdapter
}

// NewTokenRepository returns initialized Datastore
func NewTokenRepository() (repository.TokenRepository, error) {
	db, err := utils.ConnectMySQL()
	if err != nil {
		return nil, err
	}

	return &TokenRepositoryImpl{
		&SQLRepositoryAdapter{Conn: db},
	}, nil
}

func (r *TokenRepositoryImpl) mapToEntity(row *sql.Row) (*domain.Token, error) {
	m := &domain.Token{}
	err := row.Scan(&m.ID, &m.Value)
	return m, err
}

// Get return a token record matched by 'id'
func (r *TokenRepositoryImpl) Get(id int) (*domain.Token, error) {
	row, err := r.queryRow("select id, value from tokens where id=?", id)
	if err != nil {

		return nil, err
	}
	d, err := r.mapToEntity(row)
	if err != nil && err == sql.ErrNoRows {
		return nil, domain.ErrNotFoundToken
	}
	return d, err
}

// FindByValue return a token record matched by 'value'
func (r *TokenRepositoryImpl) FindByValue(value string) (*domain.Token, error) {
	row, err := r.queryRow("select id, value from tokens where value=?", value)
	if err != nil {
		return nil, err
	}
	d, err := r.mapToEntity(row)
	if err != nil && err == sql.ErrNoRows {
		return nil, domain.ErrNotFoundToken
	}
	return d, err
}

// Save saves token data to datastore
func (r *TokenRepositoryImpl) Save(m *domain.Token) (int, error) {
	_, err := r.FindByValue(m.Value)
	if err == nil {
		return 0, domain.ErrTokenAlreadyExistSameValue
	}

	stmt, err := r.Conn.Prepare("insert into tokens (value) values(?)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(m.Value)
	if err != nil {
		return 0, err
	}
	id, _ := res.LastInsertId()
	return int(id), nil
}

// Update update the value
func (r *TokenRepositoryImpl) Update(m *domain.Token) error {
	stmt, err := r.Conn.Prepare("update tokens set value=? where id=?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(m.Value, m.ID)
	return err
}

// Delete deletes record when matched id
// Returns number of deleted record and an error
func (r *TokenRepositoryImpl) Delete(id int) (bool, error) {
	stmt, err := r.Conn.Prepare("delete from tokens where id=?")
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(id)
	cnt, _ := res.RowsAffected()
	return cnt > 0, err
}
