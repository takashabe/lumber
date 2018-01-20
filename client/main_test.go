package client

import (
	"os"
	"testing"

	"github.com/takashabe/lumber/helper"
)

func TestMain(m *testing.M) {
	setup()
	os.Exit(m.Run())
}

func setup() {
	helper.SetupTables()
	os.Setenv(LumberToken, "foo")
}
