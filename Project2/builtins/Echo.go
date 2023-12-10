package builtins

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrInvalidArgEcho = errors.New("invalid argument count for Echo")
)

// Echo prints arguments to the standard output.
func Echo(args ...string) error {
	if len(args) < 1 {
		return ErrInvalidArgEcho
	}
	fmt.Println(strings.Join(args, " "))
	return nil
}
