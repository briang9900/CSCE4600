package builtins

import (
	"bytes"
	"testing"
)

func TestPwd(t *testing.T) {
	// Redirect output to a buffer
	buf := new(bytes.Buffer)
	err := Pwd(buf)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// For Pwd, no need to compare output, just ensure no error occurred
}
