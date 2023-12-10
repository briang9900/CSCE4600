package builtins

import (
	"testing"
)

func TestPwd(t *testing.T) {
	// Call ChangeDirectory without any arguments
	err := ChangeDirectory()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// For ChangeDirectory, no need to compare output, just ensure no error occurred
}
