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

	// For Ls, no need to compare output, just ensure no error occurred
	// If you want to check directory contents, you can read from the buffer
	// Example:
	// fmt.Println("Directory contents:")
	// fmt.Println(buf.String())
}
