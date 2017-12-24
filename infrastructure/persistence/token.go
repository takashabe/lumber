package persistence

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/takashabe/lumber/domain"
)

// TokenRepositoryImpl implements the TokenRepository
type TokenRepositoryImpl struct {
	*SQLRepositoryAdapter
}

func (r *TokenRepositoryImpl) mapToEntity(row *sql.Row) (*domain.Token, error) {
	m := &domain.Token{}
	err := row.Scan(&m.ID, &m.Value, &m.CreatedAt, &m.UpdatedAt)
	return m, err
}

// Get return a token record matched by 'id'
func (r *TokenRepositoryImpl) Get(id int) (*domain.Token, error) {
	row, err := r.queryRow("select * from tokens where id=?", id)
	if err != nil {
		return nil, err
	}
	return r.mapToEntity(row)
}

// FindByValue return a token record matched by 'value'
func (r *TokenRepositoryImpl) FindByValue(value string) (*domain.Token, error) {
	row, err := r.queryRow("select * from tokens where value=?", value)
	if err != nil {
		return nil, err
	}
	return r.mapToEntity(row)
}

// Save saves token data to datastore
func (r *TokenRepositoryImpl) Save(m *domain.Token) (int, error) {
	_, err := r.FindByValue(m.Value)
	if err == nil {
		return 0, errors.New("failed to save token. A record with the same value already exists")
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
