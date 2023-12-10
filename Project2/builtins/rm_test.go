package builtins

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestRm(t *testing.T) {
	testFile := "testFile"

	// Create a test file
	if _, err := os.Create(testFile); err != nil {
		t.Fatalf("Error creating test file: %s", err)
	}

	err := Rm(testFile)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Check if file was removed
	_, err = ioutil.ReadFile(testFile)
	if !os.IsNotExist(err) {
		t.Errorf("File was not removed: %s", err)
	}
}
