package builtins

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestMkdir(t *testing.T) {
	testDir := "testDir"

	// Clean up test directory after testing
	defer func() {
		if err := os.RemoveAll(testDir); err != nil {
			t.Fatalf("Error cleaning up test directory: %s", err)
		}
	}()

	err := Mkdir(testDir)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Check if directory exists
	_, err = ioutil.ReadDir(testDir)
	if err != nil {
		t.Errorf("Error reading directory: %s", err)
	}
}
