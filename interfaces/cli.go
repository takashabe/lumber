package interfaces

import (
	"fmt"
	"io"

	"github.com/takashabe/lumber/infrastructure/persistence"
	"github.com/takashabe/lumber/library/config"
)

// Exit codes. used only in Run()
const (
	ExitCodeOK = 0

	// Specific error codes. begin 10-
	ExitCodeError = 10 + iota
	ExitCodeSetupServerError
)

// CLI is the command line interface object
type CLI struct {
	OutStream io.Writer
	ErrStream io.Writer
}

// Run invokes the CLI with the given arguments
func (c *CLI) Run(args []string) int {
	conf := config.Config.Server

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

	if err := server.Run(conf.Port); err != nil {
		fmt.Fprintf(c.ErrStream, "failed from server: %v", err)
		return ExitCodeError
	}
	return ExitCodeOK
}
