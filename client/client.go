package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/pkg/errors"
)

// Constants related to environment variables
const (
	LumberServerAddress = "LUMBER_SERVER_ADDRESS"
	LumberToken         = "LUMBER_SESSION_TOKEN"
)

// Client represent a client for the lumber server
type Client struct {
	addr  string
	token string
}

// New returns initialized client
func New() (*Client, error) {
	token := os.Getenv(LumberToken)
	addr := os.Getenv(LumberServerAddress)
	if len(addr) == 0 {
		return nil, errors.New("Require settings server address in environment variable")
	}
	if addr[len(addr)-1] != '/' {
		addr = addr + "/"
	}
	if _, err := url.Parse(addr); err != nil {
		return nil, err
	}

	return &Client{
		addr:  addr,
		token: token,
	}, nil
}

// CreateEntry submit markdown file as a new entry
func (c *Client) CreateEntry(ctx context.Context, file string) (int, error) {
	if len(c.token) == 0 {
		return 0, ErrRequireToken
	}

	f, err := ioutil.ReadFile(file)
	if err != nil {
		return 0, err
	}

	type payload struct {
		Data   []byte `json:"data"`
		Status int    `json:"status"`
	}
	raw := payload{
		Data:   f,
		Status: 1, // TODO: changeable status
	}
	var buf bytes.Buffer
	err = json.NewEncoder(&buf).Encode(raw)
	if err != nil {
		return 0, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%sapi/entry?token=%s", c.addr, c.token), &buf)
	if err != nil {
		return 0, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()
	err = verifyHTTPStatusCode(res, http.StatusOK, http.StatusNoContent)
	if err != nil {
		// TODO: error handling
		b, _ := ioutil.ReadAll(res.Body)
		return 0, errors.Wrapf(err, "response: %s", b)
	}
	if res.StatusCode == http.StatusNoContent {
		return 0, nil
	}

	type response struct {
		ID int `json:"id"`
	}
	resPayload := response{}
	err = json.NewDecoder(res.Body).Decode(&resPayload)
	if err != nil {
		return 0, err
	}

	return resPayload.ID, nil
}

// Entry returns initialized Entry
func (c *Client) Entry(id int) *Entry {
	return &Entry{
		id:    id,
		addr:  c.addr,
		token: c.token,
	}
}

func verifyHTTPStatusCode(res *http.Response, codes ...int) error {
	for _, c := range codes {
		if res.StatusCode == c {
			return nil
		}
	}
	return errors.Errorf("HTTP response error: expect status codes %v, but received %d", codes, res.StatusCode)
}
