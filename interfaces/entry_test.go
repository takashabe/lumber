package interfaces

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"

	"github.com/takashabe/lumber/helper"
)

func TestGetEntry(t *testing.T) {
	ts := setupServer(t)
	defer ts.Close()

	cases := []struct {
		input      int
		expectBody []byte
		expectCode int
	}{
		{
			1,
			[]byte(`{"id":1,"title":"foo","content":"bar","status":1}`),
			http.StatusOK,
		},
		{
			0,
			[]byte(`{"reason":"failed to get entry"}`),
			http.StatusNotFound,
		},
	}
	for i, c := range cases {
		helper.LoadFixture(t, "testdata/entries.yml")
		res := sendRequest(t, "GET", fmt.Sprintf("%s/api/entry/%d", ts.URL, c.input), nil)
		defer res.Body.Close()

		if res.StatusCode != c.expectCode {
			t.Errorf("#%d: want %d, got %d", i, c.expectCode, res.StatusCode)
		}
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Fatalf("#%d: want non error, got %v", i, err)
		}
		if !reflect.DeepEqual(body, c.expectBody) {
			t.Errorf("#%d: want body %q, got %q", i, c.expectBody, body)
		}
	}
}

func TestGetIDsEntry(t *testing.T) {
	ts := setupServer(t)
	defer ts.Close()

	cases := []struct {
		fixture    string
		expectBody []byte
		expectCode int
	}{
		{"testdata/entries.yml", []byte(`{"ids":[1,2]}`), http.StatusOK},
		{"testdata/truncate_entries.sql", []byte(`{"ids":[]}`), http.StatusOK},
	}
	for i, c := range cases {
		helper.LoadFixture(t, c.fixture)
		res := sendRequest(t, "GET", fmt.Sprintf("%s/api/entries", ts.URL), nil)
		defer res.Body.Close()

		if res.StatusCode != c.expectCode {
			t.Errorf("#%d: want %d, got %d", i, c.expectCode, res.StatusCode)
		}
		act, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Fatalf("#%d: want non error, got %#v", i, err)
		}
		if !reflect.DeepEqual(act, c.expectBody) {
			t.Errorf("#%d: want %s, got %s", i, c.expectBody, act)
		}
	}
}

func TestGetTitlesEntry(t *testing.T) {
	ts := setupServer(t)
	defer ts.Close()

	cases := []struct {
		fixture    string
		start      int
		length     int
		expectBody []byte
		expectCode int
	}{
		{
			"testdata/entries.yml",
			0,
			2,
			[]byte(`{"data":[{"id":1,"title":"foo"},{"id":2,"title":"foo"}]}`),
			http.StatusOK,
		},
	}
	for i, c := range cases {
		helper.LoadFixture(t, c.fixture)
		res := sendRequest(t, "GET", fmt.Sprintf("%s/api/titles/%d/%d", ts.URL, c.start, c.length), nil)
		defer res.Body.Close()

		if res.StatusCode != c.expectCode {
			t.Errorf("#%d: want %d, got %d", i, c.expectCode, res.StatusCode)
		}
		act, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Fatalf("#%d: want non error, got %#v", i, err)
		}
		if !reflect.DeepEqual(act, c.expectBody) {
			t.Errorf("#%d: want %s, got %s", i, c.expectBody, act)
		}
	}
}

func TestPostEntry(t *testing.T) {
	ts := setupServer(t)
	defer ts.Close()

	helper.LoadFixture(t, "testdata/entries.yml")
	helper.LoadFixture(t, "testdata/tokens.yml")

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
				Status: 1,
			},
			"foo",
			http.StatusNoContent, // duplicate title
		},
		{
			postPayload{
				Data:   []byte("# title99\n\n## content"),
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
			"notfound",
			http.StatusUnauthorized,
		},
		{
			postPayload{
				Data:   []byte("# title\n\n## content"),
				Status: 1,
			},
			"",
			http.StatusUnauthorized,
		},
	}
	for i, c := range cases {
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
			http.StatusUnauthorized,
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
		{1, "", http.StatusUnauthorized},
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
