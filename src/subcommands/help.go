package subcommands

import (
	"errors"
	"fmt"
)

func BasicHelp() error {
	fmt.Printf("Semi is a tool for building or running semi programs.\n\n")
	fmt.Printf("Usage:\n\n")
	fmt.Printf("        semi [--help] [--version]\n")
	fmt.Printf("             <command> [arguments]\n\n")
	fmt.Printf("The commands are:\n\n")
	fmt.Printf("        build         compile semi program and dependencies\n")
	fmt.Printf("        help          print help information about a topic or command\n")
	fmt.Printf("        parse-tree    print the parse tree of a program\n")
	fmt.Printf("        run           compile and run semi program\n")
	fmt.Printf("        version       print Semi version\n\n")
	fmt.Printf("Use \"semi help <command>\" for more information about a command.\n\n")

	return nil
}

func CommandHelp(command string) error {
	// TODO: Add option and argument descriptions
	switch command {
	case "build":
		fmt.Printf("usage semi build [-o output] <source>\n")
	case "help":
		fmt.Printf("usage: semi help <command>\n")
	case "parse-tree":
		fmt.Printf("usage semi parse-tree [-o output] <source>\n")
	case "run":
		fmt.Printf("usage: semi run <source> <arguments>\n")
	case "version":
		fmt.Printf("usage semi version\n")
	default:
		return errors.New("unrecognized sub-command")
	}

	return nil
}
