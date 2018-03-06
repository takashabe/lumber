package main

import (
	"os"

	"github.com/takashabe/lumber/client"
)

func main() {
	c := &client.CLI{
		OutStream: os.Stdout,
		ErrStream: os.Stderr,
	}
	os.Exit(c.Run(os.Args))
}
