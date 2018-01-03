package application

import (
	"testing"

	"github.com/takashabe/lumber/domain"
	"github.com/takashabe/lumber/helper"
	"github.com/takashabe/lumber/infrastructure/persistence"
)

func TestGetToken(t *testing.T) {
	repo, err := persistence.NewTokenRepository()
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
		interactor := NewTokenInteractor(repo)
		token, err := interactor.Get(c.input)
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

func TestNewToken(t *testing.T) {
	repo, err := persistence.NewTokenRepository()
	if err != nil {
		t.Fatalf("want non error, got %#v", err)
	}

	interactor := NewTokenInteractor(repo)
	token, err := interactor.New()
	if err != nil {
		t.Fatalf("want non error, got %#v", err)
	}

	after, err := interactor.Get(token.ID)
	if err != nil {
		t.Fatalf("want non error, got %#v", err)
	}
	if token.ID != after.ID {
		t.Errorf("want id %d, got %d", after.ID, token.ID)
	}
}
