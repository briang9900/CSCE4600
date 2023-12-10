package builtins

import (
	"errors"
	"fmt"
	"io/ioutil"
)

var (
	ErrInvalidArgLs = errors.New("invalid argument count for Ls")
)

// Ls lists directory contents.
func Ls() error {
	files, err := ioutil.ReadDir(".")
	if err != nil {
		return err
	}
	for _, file := range files {
		fmt.Println(file.Name())
	}
	return nil
}
