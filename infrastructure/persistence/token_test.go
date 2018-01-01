package persistence

import (
	"testing"

	"github.com/takashabe/lumber/domain"
	"github.com/takashabe/lumber/helper"
)

func TestGetToken(t *testing.T) {
	repo, err := NewTokenRepository()
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
		{0, 0, domain.ErrNotFoundToken},
	}
	for i, c := range cases {
		token, err := repo.Get(c.input)
		if err != c.expectErr {
			t.Errorf("#%d: want error %#v, got %#v", i, c.expectErr, err)
		}
		if err != nil {
			continue
		}

		if token.ID != c.expectID {
			t.Errorf("#%d: want id %d, got %d", i, c.expectID, token.ID)
		}
	}
}

func TestFindByValueToken(t *testing.T) {
	repo, err := NewTokenRepository()
	if err != nil {
		t.Fatalf("want non error, got %#v", err)
	}
	helper.LoadFixture(t, "testdata/tokens.yml")

	cases := []struct {
		input     string
		expectID  int
		expectErr error
	}{
		{"foo", 1, nil},
		{"bar", 2, nil},
		{"", 0, domain.ErrNotFoundToken},
	}
	for i, c := range cases {
		token, err := repo.FindByValue(c.input)
		if err != c.expectErr {
			t.Errorf("#%d: want error %#v, got %#v", i, c.expectErr, err)
		}
		if err != nil {
			continue
		}

		if token.ID != c.expectID {
			t.Errorf("#%d: want id %d, got %d", i, c.expectID, token.ID)
		}
	}
}

func TestSaveToken(t *testing.T) {
	repo, err := NewTokenRepository()
	if err != nil {
		t.Fatalf("want non error, got %#v", err)
	}

	cases := []struct {
		input     *domain.Token
		expectErr error
	}{
		{
			&domain.Token{Value: "test"},
			nil,
		},
		{
			&domain.Token{Value: "foo"},
			domain.ErrTokenAlreadyExistSameValue,
		},
	}
	for i, c := range cases {
		helper.LoadFixture(t, "testdata/tokens.yml")
		id, err := repo.Save(c.input)
		if err != c.expectErr {
			t.Errorf("#%d: want error %#v, got %#v", i, c.expectErr, err)
		}

		if err != nil {
			continue
		}
		token, err := repo.Get(id)
		if err != nil {
			t.Errorf("#%d: want non error, got %#v", i, err)
		}
		if token.Value != c.input.Value {
			t.Errorf("#%d: want value %#v, got %#v", i, c.input.Value, token.Value)
		}
	}
}
