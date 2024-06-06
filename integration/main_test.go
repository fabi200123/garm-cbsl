package integration

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	StartClient()
	result := m.Run()

	os.Exit(result)
}
