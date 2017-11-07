package model

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/takashabe/lumber/config"
)

func TestNewEntry(t *testing.T) {
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
			t.Fatalf("#%d: want non error, got %v", i, err)
		}
		entry, err := NewEntry(data, EntryStatusPublic)
		if err != c.expectErr {
			t.Fatalf("#%d: want error %v, got %v", i, c.expectErr, err)
		}
		if err != nil {
			continue
		}

		if entry.title != c.expectTitle {
			t.Fatalf("#%d: want title %s, got %s", i, c.expectTitle, entry.title)
		}
		if !strings.HasPrefix(entry.content, c.expectContentHead) {
			t.Fatalf("#%d: want content prefix %s", i, c.expectContentHead)
		}
		if !strings.HasSuffix(entry.content, c.expectContentTail) {
			t.Fatalf("#%d: want content suffix %s", i, c.expectContentTail)
		}
	}
}

func TestPost(t *testing.T) {
	cases := []struct {
		inputFilePath string
	}{
		{"testdata/go-pubsub_readme.md"},
		{"testdata/minimum.md"},
	}
	for i, c := range cases {
		data, err := ioutil.ReadFile(c.inputFilePath)
		if err != nil {
			t.Fatalf("#%d: want non error, got %v", i, err)
		}
		entry, err := NewEntry(data, EntryStatusPublic)
		if err != nil {
			t.Fatalf("#%d: want non error, got %v", i, err)
		}
		err = entry.Post()
		if err != nil {
			t.Fatalf("#%d: want non error, got %v", i, err)
		}
	}
}
