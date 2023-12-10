package builtins

import (
	"testing"
)

func TestLs(t *testing.T) {
	// Call Ls without any arguments
	err := Ls()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

}
