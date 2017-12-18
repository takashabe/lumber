package application

import (
	"database/sql"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/takashabe/lumber/config"
	"github.com/takashabe/lumber/domain/repository"
	"github.com/takashabe/lumber/helper"
	"github.com/takashabe/lumber/infrastructure/persistence"
)

func getRepository(t *testing.T) repository.EntryRepository {
	// TODO: selectable mock or production repository
	r, err := persistence.NewDatastore()
	if err != nil {
		t.Fatalf("want non error, got %#v", err)
	}
	return r
}

func TestNewEntryElement(t *testing.T) {
	cases := []struct {
		inputFilePath     string
		expectTitle       string
		expectContentHead string
		expectContentTail string
		expectErr         error
	}{
		{
			"testdata/go-pubsub_readme.md",
			"go-pubsub",
			"<p><a href=",
			"</ul>",
			nil,
		},
		{
			"testdata/minimum.md",
			"title",
			"<p>content</p>",
			"<p>content</p>",
			nil,
		},
		{
			"testdata/empty.md",
			"",
			"",
			"",
			config.ErrEmptyEntry,
		},
	}
	for i, c := range cases {
		data, err := ioutil.ReadFile(c.inputFilePath)
		if err != nil {
			t.Fatalf("#%d: want non error, got %#v", i, err)
		}
		entry, err := NewEntryElement(data)
		if err != c.expectErr {
			t.Fatalf("#%d: want error %#v, got %#v", i, c.expectErr, err)
		}
		if err != nil {
			continue
		}

		if entry.Title != c.expectTitle {
			t.Fatalf("#%d: want title %s, got %s", i, c.expectTitle, entry.Title)
		}
		if !strings.HasPrefix(entry.Content, c.expectContentHead) {
			t.Fatalf("#%d: want content prefix %s", i, c.expectContentHead)
		}
		if !strings.HasSuffix(entry.Content, c.expectContentTail) {
			t.Fatalf("#%d: want content suffix %s", i, c.expectContentTail)
		}
	}
}

func TestGetEntry(t *testing.T) {
	cases := []struct {
		input     int
		expectID  int
		expectErr error
	}{
		{1, 1, nil},
		{0, 0, sql.ErrNoRows},
	}
	for i, c := range cases {
		helper.LoadFixture(t, "testdata/entries.yml")

		interactor := NewEntryInteractor(getRepository(t))
		act, err := interactor.Get(c.input)
		if err != c.expectErr {
			t.Fatalf("#%d: want error %#v, got %#v", i, c.expectErr, err)
		}
		if err != nil {
			continue
		}

		if act.ID != c.expectID {
			t.Fatalf("#%d: want %d, got %d", i, c.expectID, act.ID)
		}
	}
}

func TestPostEntry(t *testing.T) {
	cases := []struct {
		inputFilePath string
	}{
		{"testdata/go-pubsub_readme.md"},
		{"testdata/minimum.md"},
	}
	for i, c := range cases {
		data, err := ioutil.ReadFile(c.inputFilePath)
		if err != nil {
			t.Fatalf("#%d: want non error, got %#v", i, err)
		}
		element, err := NewEntryElement(data)
		if err != nil {
			t.Fatalf("#%d: want non error, got %#v", i, err)
		}
		interactor := NewEntryInteractor(getRepository(t))
		_, err = interactor.Post(element)
		if err != nil {
			t.Fatalf("#%d: want non error, got %#v", i, err)
		}
	}
}

func TestEditEntry(t *testing.T) {
	cases := []struct {
		inputID      int
		inputData    []byte
		expectGetErr error
	}{
		{1, []byte("# title\n\n## content"), nil},
		{0, []byte("# title\n\n## content"), sql.ErrNoRows},
	}
	for i, c := range cases {
		helper.LoadFixture(t, "testdata/entries.yml")

		element, err := NewEntryElement(c.inputData)
		if err != nil {
			t.Fatalf("#%d: want non error, got %#v", i, err)
		}
		interactor := NewEntryInteractor(getRepository(t))
		err = interactor.Edit(c.inputID, element)
		if err != nil {
			t.Fatalf("#%d: want non error, got %#v", i, err)
		}

		entry, err := interactor.Get(c.inputID)
		if err != c.expectGetErr {
			t.Fatalf("#%d: want error %#v, got %#v", i, c.expectGetErr, err)
		}
		if err != nil {
			continue
		}

		if entry.Title != element.Title || entry.Content != element.Content {
			t.Fatalf("#%d: want title %s and content %s, got title %s and content %s",
				i, entry.Title, entry.Content, element.Title, element.Content)
		}
	}
}

func TestDelete(t *testing.T) {
	cases := []struct {
		input int
	}{
		{1},
		{0},
	}
	for i, c := range cases {
		helper.LoadFixture(t, "testdata/entries.yml")

		interactor := NewEntryInteractor(getRepository(t))
		err := interactor.Delete(c.input)
		if err != nil {
			t.Fatalf("#%d: want non error, got %#v", i, err)
		}
		_, err = interactor.Get(c.input)
		if err != sql.ErrNoRows {
			t.Fatalf("#%d: want error sql.ErrNoRows, got %#v", i, err)
		}
	}
}
