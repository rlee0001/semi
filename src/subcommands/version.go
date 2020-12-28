package subcommands

import (
	"errors"
	"fmt"
)

func ParseVersionArguments(arguments []string) error {
	if len(arguments) != 0 {
		return errors.New("command error: unrecognized arguments")
	}

	return nil
}

func Version() error {
	// TODO: Print clang/llvm versions.
	fmt.Printf("semi version 0.0.1\n")

	return nil
}
