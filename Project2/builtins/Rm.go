package builtins

import (
	"errors"
	//"fmt"
	"os"
)

var (
	ErrInvalidArgRm = errors.New("invalid argument count for Rm")
)

// Rm removes files or directories.
func Rm(args ...string) error {
	if len(args) < 1 {
		return ErrInvalidArgRm
	}
	for _, file := range args {
		err := os.RemoveAll(file)
		if err != nil {
			return err
		}
	}
	return nil
}
