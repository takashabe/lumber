package interfaces

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/takashabe/lumber/helper"
	"github.com/takashabe/lumber/infrastructure/persistence"
)

func setupServer(t *testing.T) *httptest.Server {
	repo, err := persistence.NewDatastore()
	if err != nil {
		t.Fatalf("want non error, got %#v", err)
	}
	server := &Server{
		Entry: NewEntryHandler(repo),
	}
	return httptest.NewServer(server.Routes())
}

func sendRequest(t *testing.T, method, url string, body io.Reader) *http.Response {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		t.Fatalf("want non error, got %#v", err)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("want non error, got %#v", err)
	}
	return res
}

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
		helper.LoadFixture(t, "fixture/entries.yml")
		res := sendRequest(t, "GET", fmt.Sprintf("%s/api/entry/%d", ts.URL, c.input), nil)
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
		expect int
	}{
		{
			postPayload{
				Data:   []byte("# title\n\n## content"),
				Status: 1,
			},
			http.StatusOK,
		},
		{
			postPayload{
				Data:   []byte("# title\n\n## content"),
				Status: 99,
			},
			http.StatusNotFound,
		},
	}
	for i, c := range cases {
		helper.LoadFixture(t, "fixture/entries.yml")
		var buf bytes.Buffer
		err := json.NewEncoder(&buf).Encode(c.input)
		if err != nil {
			t.Fatalf("#%d: want non error, got %#v", i, err)
		}
		res := sendRequest(t, "POST", fmt.Sprintf("%s/api/entry/", ts.URL), &buf)
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
		expect       int
	}{
		{
			1,
			editPayload{
				Data: []byte("# title\n\n## content"),
			},
			http.StatusOK,
		},
		{
			1,
			editPayload{},
			http.StatusNotFound,
		},
		{
			0,
			editPayload{
				Data: []byte("# title\n\n## content"),
			},
			http.StatusNotFound,
		},
	}
	for i, c := range cases {
		helper.LoadFixture(t, "fixture/entries.yml")
		var buf bytes.Buffer
		err := json.NewEncoder(&buf).Encode(c.inputPayload)
		if err != nil {
			t.Fatalf("#%d: want non error, got %#v", i, err)
		}
		res := sendRequest(t, "PUT", fmt.Sprintf("%s/api/entry/%d", ts.URL, c.inputID), &buf)
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
		expect int
	}{
		{1, http.StatusOK},
		{0, http.StatusNotFound},
	}
	for i, c := range cases {
		helper.LoadFixture(t, "fixture/entries.yml")
		res := sendRequest(t, "DELETE", fmt.Sprintf("%s/api/entry/%d", ts.URL, c.input), nil)
		defer res.Body.Close()

		if res.StatusCode != c.expect {
			t.Errorf("#%d: want %d, got %d", i, c.expect, res.StatusCode)
		}
	}
}
