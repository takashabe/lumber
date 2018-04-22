package persistence

import (
	"database/sql"
	"errors"

	"github.com/takashabe/lumber/config"
	"github.com/takashabe/lumber/domain"
	"github.com/takashabe/lumber/domain/repository"
	"github.com/takashabe/lumber/infrastructure/utils"
)

// EntryRepositoryImpl implements the EntryRepository
type EntryRepositoryImpl struct {
	*SQLRepositoryAdapter
}

// NewEntryRepository returns initialized Datastore
func NewEntryRepository() (repository.EntryRepository, error) {
	db, err := utils.ConnectMySQL()
	if err != nil {
		return nil, err
	}

	return &EntryRepositoryImpl{
		&SQLRepositoryAdapter{Conn: db},
	}, nil
}

func (r *EntryRepositoryImpl) mapToEntity(row *sql.Row) (*domain.Entry, error) {
	m := &domain.Entry{}
	err := row.Scan(&m.ID, &m.Title, &m.Content, &m.Status)
	return m, err
}

// Get return a entry record matched by 'id'
func (r *EntryRepositoryImpl) Get(id int) (*domain.Entry, error) {
	row, err := r.queryRow("select id, title, content, status from entries where id=?", id)
	if err != nil {
		return nil, err
	}
	return r.mapToEntity(row)
}

// GetByTitle return a entry record matched by 'title'
func (r *EntryRepositoryImpl) GetByTitle(title string) (*domain.Entry, error) {
	row, err := r.queryRow("select id, title, content, status from entries where title=?", title)
	if err != nil {
		return nil, err
	}
	return r.mapToEntity(row)
}

// GetIDs return all entry id list
func (r *EntryRepositoryImpl) GetIDs() ([]int, error) {
	rows, err := r.query("select id from entries")
	if err != nil {
		return nil, err
	}
	ids := make([]int, 0)
	for rows.Next() {
		var i int
		err := rows.Scan(&i)
		if err != nil {
			return nil, err
		}
		ids = append(ids, i)
	}
	return ids, nil
}

// GetTitles returns entries with contain id and title
func (r *EntryRepositoryImpl) GetTitles(start, n int) ([]*domain.Entry, error) {
	if start < 0 {
		return nil, errors.New("invalid start index")
	}
	if n < 1 {
		// default
		n = 100
	}

	// NOTE: depends on id order
	rows, err := r.query("select id, title from entries where id >= ? limit ?", start, n)
	if err != nil {
		return nil, err
	}
	entries := make([]*domain.Entry, 0)
	for rows.Next() {
		e := &domain.Entry{}
		err := rows.Scan(&e.ID, &e.Title)
		if err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}
	return entries, nil
}

// Save saves entry data to datastore
func (r *EntryRepositoryImpl) Save(e *domain.Entry) (int, error) {
	sizeTitle := len(e.Title)
	sizeContent := len(e.Content)
	if sizeTitle == 0 || sizeContent == 0 {
		return 0, config.ErrEmptyEntry
	}
	if sizeTitle > config.MaxTitleBytes || sizeContent > config.MaxContentBytes {
		return 0, config.ErrEntrySizeLimitExceeded
	}

	stmt, err := r.Conn.Prepare("insert into entries (title, content, status) values(?, ?, ?)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(e.Title, e.Content, int(e.Status))
	if err != nil {
		return 0, err
	}
	id, _ := res.LastInsertId()
	return int(id), nil
}

// Edit update the title and content of the entry
func (r *EntryRepositoryImpl) Edit(e *domain.Entry) error {
	stmt, err := r.Conn.Prepare("update entries set title=?, content=? where id=?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(e.Title, e.Content, e.ID)
	return err
}

// Delete delets record when matched id
// Returns number of deleted record and an error
func (r *EntryRepositoryImpl) Delete(id int) (bool, error) {
	stmt, err := r.Conn.Prepare("delete from entries where id=?")
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(id)
	cnt, _ := res.RowsAffected()
	return cnt > 0, err
}
