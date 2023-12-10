package builtins

import (
	"testing"
)

func TestEcho(t *testing.T) {
	testCases := []struct {
		name     string
		args     []string
		expected string
	}{
		{
			name:     "Print single argument",
			args:     []string{"hello"},
			expected: "hello\n",
		},
		{
			name:     "Print multiple arguments",
			args:     []string{"hello", "world"},
			expected: "hello world\n",
		},
		// Add more test cases as needed
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call Echo with unpacked arguments using ...
			err := Echo(tc.args...)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			// Compare output with expected (you might need to capture the output for comparison)
		})
	}
}
