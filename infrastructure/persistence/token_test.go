package persistence

import (
	"database/sql"
	"testing"

	"github.com/takashabe/lumber/helper"
)

func TestGetToken(t *testing.T) {
	db, err := NewEntryRepository()
	if err != nil {
		t.Fatalf("want non error, got %#v", err)
	}
	helper.LoadFixture(t, "testdata/tokens.yml")

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
