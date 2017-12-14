package interfaces

import (
	"os"
	"testing"

	"github.com/takashabe/lumber/helper"
)

func TestMain(m *testing.M) {
	helper.SetupTables()
	os.Exit(m.Run())
}
