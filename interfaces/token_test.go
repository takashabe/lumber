package interfaces

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/takashabe/lumber/helper"
)

func TestGetToken(t *testing.T) {
	ts := setupServer(t)
	defer ts.Close()

	cases := []struct {
		input  int
		expect int
	}{
		{1, http.StatusOK},
		{0, http.StatusNotFound},
	}
	for i, c := range cases {
		helper.LoadFixture(t, "testdata/tokens.yml")
		res := sendRequest(t, "GET", fmt.Sprintf("%s/api/token/%d", ts.URL, c.input), nil)
		defer res.Body.Close()

		if res.StatusCode != c.expect {
			t.Errorf("#%d: want %d, got %d", i, c.expect, res.StatusCode)
		}
	}
}

func TestFindByValueToken(t *testing.T) {
	ts := setupServer(t)
	defer ts.Close()

	cases := []struct {
		input  string
		expect int
	}{
		{"foo", http.StatusOK},
		{"", http.StatusNotFound},
	}
	for i, c := range cases {
		helper.LoadFixture(t, "testdata/tokens.yml")
		res := sendRequest(t, "GET", fmt.Sprintf("%s/api/token/value/%s", ts.URL, c.input), nil)
		defer res.Body.Close()

		if res.StatusCode != c.expect {
			t.Errorf("#%d: want %d, got %d", i, c.expect, res.StatusCode)
		}
	}
}
