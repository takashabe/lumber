package persistence

import (
	"bytes"
	"database/sql"
	"reflect"
	"testing"

	"github.com/takashabe/lumber/config"
	"github.com/takashabe/lumber/domain"
	"github.com/takashabe/lumber/helper"
)

func TestGetEntry(t *testing.T) {
	db, err := NewEntryRepository()
	if err != nil {
		t.Fatalf("want non error, got %#v", err)
	}
	helper.LoadFixture(t, "testdata/entries.yml")

	cases := []struct {
		input     int
		expectID  int
		expectErr error
	}{
		{1, 1, nil},
		{0, 0, sql.ErrNoRows},
	}
	for i, c := range cases {
		model, err := db.Get(c.input)
		if err != c.expectErr {
			t.Errorf("#%d: want error %#v, got %#v", i, c.expectErr, err)
		}

		if model.ID != c.expectID {
			t.Errorf("#%d: want id %d, got %d", i, c.expectID, model.ID)
		}
	}
}

func TestGetIDsEntry(t *testing.T) {
	db, err := NewEntryRepository()
	if err != nil {
		t.Fatalf("want non error, got %#v", err)
	}

	cases := []struct {
		fixture   string
		expectIDs []int
	}{
		{
			"testdata/entries.yml",
			[]int{1, 2},
		},
		{
			"testdata/delete_entries.sql",
			[]int{},
		},
	}
	for i, c := range cases {
		helper.LoadFixture(t, c.fixture)
		ids, err := db.GetIDs()
		if err != nil {
			t.Errorf("#%d: want non error, got %#v", i, err)
		}
		if !reflect.DeepEqual(ids, c.expectIDs) {
			t.Errorf("#%d: want ids %#v, got %#v", i, c.expectIDs, ids)
		}
	}
}

func TestGetTitlesEntry(t *testing.T) {
	db, err := NewEntryRepository()
	if err != nil {
		t.Fatalf("want non error, got %#v", err)
	}

	cases := []struct {
		fixture string
		start   int
		length  int
		expect  []*domain.Entry
	}{
		{
			"testdata/entries.yml",
			0,
			2,
			[]*domain.Entry{
				&domain.Entry{ID: 1, Title: "foo"},
				&domain.Entry{ID: 2, Title: "foo"},
			},
		},
		{
			"testdata/entries.yml",
			2,
			2,
			[]*domain.Entry{
				&domain.Entry{ID: 2, Title: "foo"},
			},
		},
		{
			"testdata/entries.yml",
			0,
			0,
			[]*domain.Entry{
				&domain.Entry{ID: 1, Title: "foo"},
				&domain.Entry{ID: 2, Title: "foo"},
			},
		},
		{
			"testdata/delete_entries.sql",
			0,
			2,
			[]*domain.Entry{},
		},
	}
	for i, c := range cases {
		helper.LoadFixture(t, c.fixture)
		es, err := db.GetTitles(c.start, c.length)
		if err != nil {
			t.Errorf("#%d: want non error, got %#v", i, err)
		}
		if !reflect.DeepEqual(es, c.expect) {
			t.Errorf("#%d: want %q, got %q", i, c.expect, es)
		}
	}
}

func TestSaveEntry(t *testing.T) {
	db, err := NewEntryRepository()
	if err != nil {
		t.Fatalf("want non error, got %#v", err)
	}

	stringWithBytes := func(size int) string {
		buf := bytes.NewBuffer(make([]byte, 0, size))
		for i := 0; i < size; i++ {
			if err := buf.WriteByte('a'); err != nil {
				t.Fatalf("want non error, got %#v", err)
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
		entity := &domain.Entry{
			Title:   c.inputTitle,
			Content: c.inputContent,
			Status:  domain.EntryStatus(c.inputStatus),
		}
		_, err := db.Save(entity)
		if err != c.expectErr {
			t.Errorf("#%d: want error %#v, got %#v", i, c.expectErr, err)
		}
	}
}

func TestEditEntry(t *testing.T) {
	db, err := NewEntryRepository()
	if err != nil {
		t.Fatalf("want non error, got %#v", err)
	}

	cases := []struct {
		inputID      int
		inputTitle   string
		inputContent string
	}{
		{
			1,
			"edit_title",
			"edit_content",
		},
	}
	for i, c := range cases {
		helper.LoadFixture(t, "testdata/entries.yml")

		entity := &domain.Entry{
			ID:      c.inputID,
			Title:   c.inputTitle,
			Content: c.inputContent,
		}
		err = db.Edit(entity)
		if err != nil {
			t.Fatalf("#%d: want non error, got %#v", i, err)
		}
		e, err := db.Get(c.inputID)
		if err != nil {
			t.Fatalf("#%d: want non error, got %#v", i, err)
		}
		if e.Title != c.inputTitle || e.Content != c.inputContent {
			t.Errorf("#%d: want title %s and content %s, but title %s and content %s",
				i, e.Title, e.Content, c.inputTitle, c.inputContent)
		}
	}
}

func TestDeleteEntry(t *testing.T) {
	db, err := NewEntryRepository()
	if err != nil {
		t.Fatalf("want non error, got %#v", err)
	}

	cases := []struct {
		input  int
		expect bool
	}{
		{1, true},
		{0, false},
	}
	for i, c := range cases {
		helper.LoadFixture(t, "testdata/entries.yml")

		flag, err := db.Delete(c.input)
		if err != nil {
			t.Fatalf("want non error, got %#v", err)
		}
		if flag != c.expect {
			t.Errorf("#%d: want error %#v, got %#v", i, c.expect, err)
		}
	}
}
