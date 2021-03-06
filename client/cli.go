package client

import (
	"context"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

var (
	defaultAddr = ""
)

// Exit codes. used only in Run()
const (
	ExitCodeOK = 0

	// Specific error codes. begin 10-
	ExitCodeError = 10 + iota
	ExitCodeParseError
	ExitCodeInvalidArgsError
	ExitCodeSetupClientError
	ExitCodeApplyCommandError
	ExitCodeNotFoundCommandError
)

type param struct {
	addr  string
	file  string
	dir   string
	id    int
	token string
}

// CLI is the command line interface object
type CLI struct {
	OutStream io.Writer
	ErrStream io.Writer

	client *Client
}

// Run invokes the CLI with the given arguments
func (c *CLI) Run(args []string) int {
	param := &param{}
	err := c.parseArgs(args[2:], param)
	if err != nil {
		fmt.Fprintf(c.ErrStream, "args parse error: %v\n", err)
		return ExitCodeParseError
	}

	if param.addr != "" {
		os.Setenv(LumberServerAddress, param.addr)
	}
	if param.token != "" {
		os.Setenv(LumberToken, param.token)
	}
	c.client, err = New()
	if err != nil {
		fmt.Fprintf(c.ErrStream, "failed to initialized client: %v\n", err)
		return ExitCodeSetupClientError
	}

	ctx := context.Background()
	for _, cmd := range c.commands() {
		if cmd.name == args[1] {
			if err := cmd.fn(ctx, param); err != nil {
				fmt.Fprintln(c.ErrStream, err)
				return ExitCodeApplyCommandError
			}
			return ExitCodeOK
		}
	}

	return ExitCodeNotFoundCommandError
}

func (c *CLI) parseArgs(args []string, p *param) error {
	flags := flag.NewFlagSet("param", flag.ContinueOnError)
	flags.SetOutput(c.ErrStream)

	flags.StringVar(&p.addr, "addr", defaultAddr, "Lumber server address.")
	flags.StringVar(&p.file, "file", "", "Post or Edit an entry file")
	flags.StringVar(&p.dir, "dir", "", "Post an entries in the directory")
	flags.IntVar(&p.id, "id", 0, "Specific ID of an entry")
	flags.StringVar(&p.token, "token", "", "Server token")

	err := flags.Parse(args)
	if err != nil {
		return errors.Wrapf(err, "failed to parsed args")
	}
	return nil
}

type command struct {
	name string
	desc string
	fn   func(ctx context.Context, p *param) error
}

func (c *CLI) commands() []command {
	return []command{
		{
			"post",
			"post entry",
			c.doPostEntry,
		},
		{
			"post-dir",
			"post an entries in the directory",
			c.doPostEntryWithDir,
		},
		{
			"edit",
			"edit the entry",
			c.doEditEntry,
		},
	}
}

func (c *CLI) doPostEntry(ctx context.Context, p *param) error {
	id, err := c.client.CreateEntry(ctx, p.file)
	if err != nil {
		return errors.Wrap(err, "failed post entry")
	}
	fmt.Fprintf(c.OutStream, "succeed post entry. id=%d\n", id)
	return nil
}

func (c *CLI) doPostEntryWithDir(ctx context.Context, p *param) error {
	f, err := os.Stat(p.dir)
	if err != nil {
		return err
	}
	if !f.IsDir() {
		return errors.Errorf("invalid args: %s doesn't directory", p.dir)
	}
	return c.postEntryRecursiveDir(ctx, p)
}

func (c *CLI) postEntryRecursiveDir(ctx context.Context, p *param) error {
	fs, err := ioutil.ReadDir(p.dir)
	if err != nil {
		return errors.Wrapf(err, "failed to read directory")
	}

	ids := []int{}
	for _, f := range fs {
		if f.IsDir() {
			c.postEntryRecursiveDir(ctx, p)
		}

		path := filepath.Join(p.dir, f.Name())
		id, err := c.client.CreateEntry(ctx, path)
		if err != nil {
			return err
		}
		ids = append(ids, id)
	}
	return nil
}

func (c *CLI) doEditEntry(ctx context.Context, p *param) error {
	e := c.client.Entry(p.id)
	err := e.Edit(ctx, p.file)
	if err != nil {
		return errors.Wrap(err, "failed to edit an entry")
	}
	return nil
}
