package docs

import (
	"os"
	"testing"
)

// Coverage purpose only xD
func TestMain(m *testing.M) {
	exitVal := m.Run()
	os.Exit(exitVal)
}
