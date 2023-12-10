package builtins

import (
	"errors"
	"fmt"
	"os"
)

var (
	ErrInvalidArgMkdir = errors.New("invalid argument count for Mkdir")
)

// Mkdir creates a new directory.
func Mkdir(args ...string) error {
	if len(args) < 1 {
		return ErrInvalidArgMkdir
	}
	for _, dir := range args {
		err := os.Mkdir(dir, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}
