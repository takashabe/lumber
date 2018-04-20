package application

import (
	"database/sql"
	"io/ioutil"
	"reflect"
	"strings"
	"testing"

	"github.com/pkg/errors"
	"github.com/takashabe/lumber/config"
	"github.com/takashabe/lumber/domain"
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

		interactor := NewEntryInteractor(getEntryRepository(t))
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

		interactor := NewEntryInteractor(getEntryRepository(t))
		act, err := interactor.GetIDs()
		if err != nil {
			t.Fatalf("#%d: want non error, got %#v", i, err)
		}
		if !reflect.DeepEqual(act, c.expectIDs) {
			t.Fatalf("#%d: want %d, got %d", i, c.expectIDs, act)
		}
	}
}

func TestGetTitlesEntry(t *testing.T) {
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
	}
	for i, c := range cases {
		helper.LoadFixture(t, c.fixture)

		interactor := NewEntryInteractor(getEntryRepository(t))
		act, err := interactor.GetTitles(c.start, c.length)
		if err != nil {
			t.Fatalf("#%d: want non error, got %#v", i, err)
		}
		if !reflect.DeepEqual(act, c.expect) {
			t.Fatalf("#%d: want %q, got %q", i, c.expect, act)
		}
	}
}

func TestPostEntry(t *testing.T) {
	cases := []struct {
		inputFilePath string
		expectErr     error
	}{
		{"testdata/go-pubsub_readme.md", nil},
		{"testdata/minimum.md", nil},
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
		interactor := NewEntryInteractor(getEntryRepository(t))
		_, err = interactor.Post(element)
		if errors.Cause(err) != c.expectErr {
			t.Fatalf("#%d: want %#v, got %#v", i, c.expectErr, err)
		}
	}
}

func TestPostWithPrivateTitle(t *testing.T) {
	cases := []struct {
		inputFilePath string
		expectEntry   *domain.Entry
	}{
		{
			"testdata/minimum.md",
			&domain.Entry{
				Title:   "title",
				Content: "<p>content</p>",
				Status:  0,
			},
		},
		{
			"testdata/wip.md",
			&domain.Entry{
				Title:   "[wip] title",
				Content: "<p>content</p>",
				Status:  1,
			},
		},
	}
	for i, c := range cases {
		data, _ := ioutil.ReadFile(c.inputFilePath)
		element, _ := NewEntryElement(data)
		interactor := NewEntryInteractor(getEntryRepository(t))
		id, err := interactor.Post(element)
		if err != nil {
			t.Fatalf("#%d: want non error, got %#v", i, err)
		}

		c.expectEntry.ID = id
		actual, err := interactor.Get(id)
		if err != nil {
			t.Fatalf("#%d: want non error, got %#v", i, err)
		}
		if !reflect.DeepEqual(actual, c.expectEntry) {
			t.Errorf("#%d: want %#v, got %#v", i, c.expectEntry, actual)
		}
	}
}

func TestEditEntry(t *testing.T) {
	cases := []struct {
		inputID   int
		inputData []byte
		err       error
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
		interactor := NewEntryInteractor(getEntryRepository(t))
		err = interactor.Edit(c.inputID, element)
		if err != nil {
			continue
		}

		entry, err := interactor.Get(c.inputID)
		if errors.Cause(err) != c.err {
			t.Fatalf("#%d: want error %#v, got %#v", i, c.err, err)
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
		expectErr error
	}{
		{1, nil},
		{0, nil},
	}
	for i, c := range cases {
		helper.LoadFixture(t, "testdata/entries.yml")

		interactor := NewEntryInteractor(getEntryRepository(t))
		err := interactor.Delete(c.input)
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
