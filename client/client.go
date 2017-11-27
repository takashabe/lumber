package client

import (
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
	if addr[len(addr)-1] != "/" {
		addr = addr + "/"
	}
	if _, err := url.Parse(addr); err != nil {
		return nil, err
	}

	return &Client{
		addr:  addr,
		token: os.Getenv(LumberSessionToken),
	}
}
