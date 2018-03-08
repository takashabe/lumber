package interfaces

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/takashabe/lumber/helper"
)

func TestGetEntry(t *testing.T) {
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
		helper.LoadFixture(t, "testdata/entries.yml")
		res := sendRequest(t, "GET", fmt.Sprintf("%s/api/entry/%d", ts.URL, c.input), nil)
		defer res.Body.Close()

		if res.StatusCode != c.expect {
			t.Errorf("#%d: want %d, got %d", i, c.expect, res.StatusCode)
		}
	}
}

func TestGetIDsEntry(t *testing.T) {
	ts := setupServer(t)
	defer ts.Close()

	cases := []struct {
		fixture string
		expect  int
	}{
		{"testdata/entries.yml", http.StatusOK},
	}
	for i, c := range cases {
		helper.LoadFixture(t, c.fixture)
		res := sendRequest(t, "GET", fmt.Sprintf("%s/api/entry/list", ts.URL), nil)
		defer res.Body.Close()

		if res.StatusCode != c.expect {
			t.Errorf("#%d: want %d, got %d", i, c.expect, res.StatusCode)
		}
	}
}

func TestPostEntry(t *testing.T) {
	ts := setupServer(t)
	defer ts.Close()

	type postPayload struct {
		Data   []byte `json:"data"`
		Status int    `json:"status"`
	}
	cases := []struct {
		input  postPayload
		token  string
		expect int
	}{
		{
			postPayload{
				Data:   []byte("# title\n\n## content"),
				Status: 1,
			},
			"foo",
			http.StatusOK,
		},
		{
			postPayload{
				Data:   []byte("# title\n\n## content"),
				Status: 99,
			},
			"foo",
			http.StatusNotFound,
		},
		{
			postPayload{
				Data:   []byte("# title\n\n## content"),
				Status: 1,
			},
			"",
			http.StatusBadRequest,
		},
	}
	for i, c := range cases {
		helper.LoadFixture(t, "testdata/entries.yml")
		helper.LoadFixture(t, "testdata/tokens.yml")
		var buf bytes.Buffer
		err := json.NewEncoder(&buf).Encode(c.input)
		if err != nil {
			t.Fatalf("#%d: want non error, got %#v", i, err)
		}
		res := sendRequest(t, "POST", fmt.Sprintf("%s/api/entry/?token=%s", ts.URL, c.token), &buf)
		defer res.Body.Close()

		if res.StatusCode != c.expect {
			t.Errorf("#%d: want %d, got %d", i, c.expect, res.StatusCode)
		}
	}
}

func TestEditEntry(t *testing.T) {
	ts := setupServer(t)
	defer ts.Close()

	type editPayload struct {
		Data []byte `json:"data"`
	}
	cases := []struct {
		inputID      int
		inputPayload editPayload
		token        string
		expect       int
	}{
		{
			1,
			editPayload{
				Data: []byte("# title\n\n## content"),
			},
			"foo",
			http.StatusOK,
		},
		{
			1,
			editPayload{},
			"foo",
			http.StatusNotFound,
		},
		{
			0,
			editPayload{
				Data: []byte("# title\n\n## content"),
			},
			"foo",
			http.StatusNotFound,
		},
		{
			1,
			editPayload{
				Data: []byte("# title\n\n## content"),
			},
			"",
			http.StatusBadRequest,
		},
	}
	for i, c := range cases {
		helper.LoadFixture(t, "testdata/entries.yml")
		var buf bytes.Buffer
		err := json.NewEncoder(&buf).Encode(c.inputPayload)
		if err != nil {
			t.Fatalf("#%d: want non error, got %#v", i, err)
		}
		res := sendRequest(t, "PUT", fmt.Sprintf("%s/api/entry/%d?token=%s", ts.URL, c.inputID, c.token), &buf)
		defer res.Body.Close()

		if res.StatusCode != c.expect {
			t.Errorf("#%d: want %d, got %d", i, c.expect, res.StatusCode)
		}
	}
}

func TestDeleteEntry(t *testing.T) {
	ts := setupServer(t)
	defer ts.Close()

	cases := []struct {
		input  int
		token  string
		expect int
	}{
		{1, "foo", http.StatusOK},
		{0, "foo", http.StatusNotFound},
		{1, "", http.StatusBadRequest},
	}
	for i, c := range cases {
		helper.LoadFixture(t, "testdata/entries.yml")
		res := sendRequest(t, "DELETE", fmt.Sprintf("%s/api/entry/%d?token=%s", ts.URL, c.input, c.token), nil)
		defer res.Body.Close()

		if res.StatusCode != c.expect {
			t.Errorf("#%d: want %d, got %d", i, c.expect, res.StatusCode)
		}
	}
}
