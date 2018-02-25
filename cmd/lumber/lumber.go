package main

import (
	"os"

	"github.com/takashabe/lumber/interfaces"
)

func main() {
	c := &interfaces.CLI{
		OutStream: os.Stdout,
		ErrStream: os.Stderr,
	}
	os.Exit(c.Run(os.Args))
}
