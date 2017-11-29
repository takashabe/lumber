package client

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/pkg/errors"
)

// Constants related to environment variables
const (
	LumberServerAddress = "lumber_server_address"
	LumberSessionToken  = "lumber_session_token"
)

// Client represent a client for the lumber server
type Client struct {
	addr  string
	token string
}

// New returns initialized client
func New() (*Client, error) {
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
		token: os.Getenv(LumberSessionToken),
	}, nil
}

// CreateEntry submit markdown file as a new entry
func (c *Client) CreateEntry(ctx context.Context, file string) (int, error) {
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

	req, err := http.NewRequest("POST", c.addr+"api/entry", &buf)
	if err != nil {
		return 0, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()
	err = verifyHTTPStatusCode(http.StatusOK, res)
	if err != nil {
		return 0, err
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

func verifyHTTPStatusCode(expect int, res *http.Response) error {
	if c := res.StatusCode; c != expect {
		return errors.Errorf("HTTP response error: expecte status code %d, but received %d", expect, c)
	}
	return nil
}