package interfaces

import (
	"flag"
	"fmt"
	"io"

	"github.com/pkg/errors"
	"github.com/takashabe/lumber/infrastructure/persistence"
)

// default parameters
const (
	defaultPort = 8080
)

// Exit codes. used only in Run()
const (
	ExitCodeOK = 0

	// Specific error codes. begin 10-
	ExitCodeError = 10 + iota
	ExitCodeParseError
	ExitCodeInvalidArgsError
	ExitCodeSetupServerError
)

var (
	// ErrParseFailed is failed to cli args parse
	ErrParseFailed = errors.New("failed to parse args")
)

type param struct {
	port int
}

// CLI is the command line interface object
type CLI struct {
	OutStream io.Writer
	ErrStream io.Writer
}

// Run invokes the CLI with the given arguments
func (c *CLI) Run(args []string) int {
	param := &param{}
	err := c.parseArgs(args[1:], param)
	if err != nil {
		fmt.Fprintf(c.ErrStream, "args parse error: %v", err)
		return ExitCodeParseError
	}

	entryRepository, err := persistence.NewEntryRepository()
	tokenRepository, err := persistence.NewTokenRepository()
	if err != nil {
		fmt.Fprintf(c.ErrStream, "failed to initialized persistence repository: %v", err)
		return ExitCodeSetupServerError
	}
	server := Server{
		Entry: NewEntryHandler(
			entryRepository,
			tokenRepository,
		),
		Token: NewTokenHandler(
			tokenRepository,
		),
	}

	if err := server.Run(param.port); err != nil {
		fmt.Fprintf(c.ErrStream, "failed from server: %v", err)
		return ExitCodeError
	}
	return ExitCodeOK
}

func (c *CLI) parseArgs(args []string, p *param) error {
	flags := flag.NewFlagSet("param", flag.ContinueOnError)
	flags.SetOutput(c.ErrStream)

	flags.IntVar(&p.port, "port", defaultPort, "Running port. require unused port.")

	err := flags.Parse(args)
	if err != nil {
		return errors.Wrapf(ErrParseFailed, err.Error())
	}
	return nil
}
