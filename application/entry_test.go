package application

import (
	"database/sql"
	"io/ioutil"
	"reflect"
	"strings"
	"testing"

	"github.com/pkg/errors"
	"github.com/takashabe/lumber/config"
	"github.com/takashabe/lumber/helper"
)

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

		interactor := NewEntryInteractor(getEntryRepository(t), getTokenRepository(t))
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

func TestGetIDsEntry(t *testing.T) {
	cases := []struct {
		fixture   string
		expectIDs []int
	}{
		{
			"testdata/entries.yml",
			[]int{1, 2},
		},
	}
	for i, c := range cases {
		helper.LoadFixture(t, c.fixture)

		interactor := NewEntryInteractor(getEntryRepository(t), getTokenRepository(t))
		act, err := interactor.GetIDs()
		if err != nil {
			t.Fatalf("#%d: want non error, got %#v", i, err)
		}
		if !reflect.DeepEqual(act, c.expectIDs) {
			t.Fatalf("#%d: want %d, got %d", i, c.expectIDs, act)
		}
	}
}

func TestPostEntry(t *testing.T) {
	cases := []struct {
		inputFilePath string
		token         string
		expectErr     error
	}{
		{"testdata/go-pubsub_readme.md", "foo", nil},
		{"testdata/minimum.md", "foo", nil},
		{"testdata/minimum.md", "", config.ErrInsufficientPrivileges},
	}
	for i, c := range cases {
		helper.LoadFixture(t, "testdata/tokens.yml")

		data, err := ioutil.ReadFile(c.inputFilePath)
		if err != nil {
			t.Fatalf("#%d: want non error, got %#v", i, err)
		}
		element, err := NewEntryElement(data)
		if err != nil {
			t.Fatalf("#%d: want non error, got %#v", i, err)
		}
		interactor := NewEntryInteractor(getEntryRepository(t), getTokenRepository(t))
		_, err = interactor.Post(element, c.token)
		if errors.Cause(err) != c.expectErr {
			t.Fatalf("#%d: want %#v, got %#v", i, c.expectErr, err)
		}
	}
}

func TestEditEntry(t *testing.T) {
	cases := []struct {
		inputID       int
		inputData     []byte
		token         string
		expectEditErr error
		expectGetErr  error
	}{
		{1, []byte("# title\n\n## content"), "foo", nil, nil},
		{0, []byte("# title\n\n## content"), "foo", nil, sql.ErrNoRows},
		{1, []byte("# title\n\n## content"), "", config.ErrInsufficientPrivileges, nil},
	}
	for i, c := range cases {
		helper.LoadFixture(t, "testdata/entries.yml")
		helper.LoadFixture(t, "testdata/tokens.yml")

		element, err := NewEntryElement(c.inputData)
		if err != nil {
			t.Fatalf("#%d: want non error, got %#v", i, err)
		}
		interactor := NewEntryInteractor(getEntryRepository(t), getTokenRepository(t))
		err = interactor.Edit(c.inputID, element, c.token)
		if errors.Cause(err) != c.expectEditErr {
			t.Fatalf("#%d: want error %#v, got %#v", i, c.expectEditErr, err)
		}
		if err != nil {
			continue
		}

		entry, err := interactor.Get(c.inputID)
		if errors.Cause(err) != c.expectGetErr {
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
		input     int
		token     string
		expectErr error
	}{
		{1, "foo", nil},
		{0, "foo", nil},
		{1, "", config.ErrInsufficientPrivileges},
	}
	for i, c := range cases {
		helper.LoadFixture(t, "testdata/entries.yml")
		helper.LoadFixture(t, "testdata/tokens.yml")

		interactor := NewEntryInteractor(getEntryRepository(t), getTokenRepository(t))
		err := interactor.Delete(c.input, c.token)
		if errors.Cause(err) != c.expectErr {
			t.Fatalf("#%d: want %#v, got %#v", i, c.expectErr, err)
		}
		if err != nil {
			continue
		}

		_, err = interactor.Get(c.input)
		if err != sql.ErrNoRows {
			t.Fatalf("#%d: want error sql.ErrNoRows, got %#v", i, err)
		}
	}
}
