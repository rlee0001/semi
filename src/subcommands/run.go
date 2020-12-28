package subcommands

import (
	"errors"
)

func ParseRunArguments(arguments []string) (string, error) {
	if len(arguments) != 1 {
		return "", errors.New("command error: <source> argument is required")
	}

	return arguments[0], nil
}

func Run(semiFilePath string) error {
	// TODO: Parse and generate ll, run "lli" command.
	_ = semiFilePath

	return nil
}
