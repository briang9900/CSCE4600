package builtins

import (
	"errors"
	"fmt"
	"os"
)

var (
	ErrInvalidArgPwd = errors.New("invalid argument count for Pwd")
)

// Pwd displays the current working directory.
func Pwd() error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	fmt.Println(wd)
	return nil
}
