package datastore

import (
	"bytes"
	"database/sql"
	"testing"

	"github.com/takashabe/go-fixture"
	"github.com/takashabe/lumber/config"
)

func TestFindEntryByID(t *testing.T) {
	db, err := NewDatastore()
	if err != nil {
		t.Fatalf("want non error, got %v", err)
	}
	f := fixture.NewFixture(db.Conn, "mysql")
	err = f.Load("fixture/entries.yml")
	if err != nil {
		t.Fatalf("want non error, got %v", err)
	}

	cases := []struct {
		input     int
		expectID  int
		expectErr error
	}{
		{1, 1, nil},
		{0, 0, sql.ErrNoRows},
	}
	for i, c := range cases {
		model, err := db.FindEntryByID(c.input)
		if err != c.expectErr {
			t.Errorf("#%d: want error %v, got %v", i, c.expectErr, err)
		}

		if model.ID != c.expectID {
			t.Errorf("#%d: want id %d, got %d", i, c.expectID, model.ID)
		}
	}
}

func TestSaveEntry(t *testing.T) {
	db, err := NewDatastore()
	if err != nil {
		t.Fatalf("want non error, got %v", err)
	}

	stringWithBytes := func(size int) string {
		buf := bytes.NewBuffer(make([]byte, 0, size))
		for i := 0; i < size; i++ {
			if err := buf.WriteByte('a'); err != nil {
				t.Fatalf("want non error, got %v", err)
			}
		}
		return buf.String()
	}

	cases := []struct {
		inputTitle   string
		inputContent string
		inputStatus  int
		expectErr    error
	}{
		{
			stringWithBytes(1 << 8),
			stringWithBytes(1<<16 - 1),
			0,
			nil,
		},
		{
			stringWithBytes(1 << 8),
			stringWithBytes(1 << 16),
			0,
			config.ErrEntrySizeLimitExceeded,
		},
		{
			"",
			"foo",
			0,
			config.ErrEmptyEntry,
		},
	}
	for i, c := range cases {
		err := db.SaveEntry(c.inputTitle, c.inputContent, c.inputStatus)
		if err != c.expectErr {
			t.Errorf("#%d: want error %v, got %v", i, c.expectErr, err)
		}
	}
}

func TestDeleteEntry(t *testing.T) {
	db, err := NewDatastore()
	if err != nil {
		t.Fatalf("want non error, got %v", err)
	}

	cases := []struct {
		input  int
		expect bool
	}{
		{1, true},
		{0, false},
	}
	for i, c := range cases {
		f := fixture.NewFixture(db.Conn, "mysql")
		err = f.Load("fixture/entries.yml")
		if err != nil {
			t.Fatalf("want non error, got %v", err)
		}
		flag, err := db.DeleteEntry(c.input)
		if err != nil {
			t.Fatalf("want non error, got %v", err)
		}
		if flag != c.expect {
			t.Errorf("#%d: want error %v, got %v", i, c.expect, err)
		}
	}
}
