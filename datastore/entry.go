package datastore

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
func (d *Datastore) SaveEntry(title, content string, status int) error {
	sizeTitle := len(title)
	sizeContent := len(content)
	if sizeTitle == 0 || sizeContent == 0 {
		return config.ErrEmptyEntry
	}
	if sizeTitle > config.MaxTitleBytes || sizeContent > config.MaxContentBytes {
		return config.ErrEntrySizeLimitExceeded
	}

	stmt, err := d.Conn.Prepare("insert into entries (title, content, status) values(?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(title, content, status)
	return err
}
