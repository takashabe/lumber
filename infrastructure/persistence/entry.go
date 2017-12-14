package persistence

import (
	"database/sql"
	"time"

	"github.com/takashabe/lumber/config"
)

// EntriesModel represent entries table
type EntriesModel struct {
	ID        int
	Title     string
	Content   string
	Status    int
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (m *EntriesModel) mapRow(r *sql.Row) error {
	return r.Scan(&m.ID, &m.Title, &m.Content, &m.Status, &m.CreatedAt, &m.UpdatedAt)
}

// FindEntryByID return a entires record matched by 'id'
func (d *Datastore) FindEntryByID(id int) (*EntriesModel, error) {
	row, err := d.queryRow("select * from entries where id=?", id)
	if err != nil {
		return nil, err
	}

	model := &EntriesModel{}
	err = model.mapRow(row)
	return model, err
}

// SaveEntry saves entry data to datastore
func (d *Datastore) SaveEntry(title, content string, status int) (int, error) {
	sizeTitle := len(title)
	sizeContent := len(content)
	if sizeTitle == 0 || sizeContent == 0 {
		return 0, config.ErrEmptyEntry
	}
	if sizeTitle > config.MaxTitleBytes || sizeContent > config.MaxContentBytes {
		return 0, config.ErrEntrySizeLimitExceeded
	}

	stmt, err := d.Conn.Prepare("insert into entries (title, content, status) values(?, ?, ?)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(title, content, status)
	if err != nil {
		return 0, err
	}
	id, _ := res.LastInsertId()
	return int(id), nil
}

// EditEntry changes the title and content of the entry
func (d *Datastore) EditEntry(id int, title, content string) error {
	stmt, err := d.Conn.Prepare("update entries set title=?, content=? where id=?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(title, content, id)
	return err
}

// DeleteEntry delets record when matched id, and returns number of deleted record and an error
func (d *Datastore) DeleteEntry(id int) (bool, error) {
	stmt, err := d.Conn.Prepare("delete from entries where id=?")
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(id)
	cnt, _ := res.RowsAffected()
	return cnt > 0, err
}
