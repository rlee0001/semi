package main

import (
	"fmt"
	"os"

	"semi/src/subcommands"
)

func main() {
	subCommand, subCommandArguments := parseArguments(os.Args[1:])

	switch subCommand {
	case "build":
		semiFilePath, binaryFileName, err := subcommands.ParseBuildArguments(subCommandArguments)
		if err != nil {
			fmt.Printf("Error: %s", err.Error())
			os.Exit(1)
		}

		err = subcommands.Build(semiFilePath, binaryFileName)
		if err != nil {
			fmt.Printf("Error: %s", err.Error())
			os.Exit(1)
		}
	case "help":
		if len(subCommandArguments) != 1 {
			err := subcommands.BasicHelp()
			if err != nil {
				fmt.Printf("Error: %s", err.Error())
				os.Exit(1)
			}
		} else {
			err := subcommands.CommandHelp(subCommandArguments[0])
			if err != nil {
				fmt.Printf("Error: %s", err.Error())
				os.Exit(1)
			}
		}
	case "parse-tree":
		semiFilePath, outputFilePath, err := subcommands.ParseParseTreeArguments(subCommandArguments)
		if err != nil {
			fmt.Printf("Error: %s", err.Error())
			os.Exit(1)
		}

		err = subcommands.ParseTree(semiFilePath, outputFilePath)
		if err != nil {
			fmt.Printf("Error: %s", err.Error())
			os.Exit(1)
		}
	case "run":
		semiFilePath, err := subcommands.ParseRunArguments(subCommandArguments)
		if err != nil {
			fmt.Printf("Error: %s", err.Error())
			os.Exit(1)
		}

		err = subcommands.Run(semiFilePath)
		if err != nil {
			fmt.Printf("Error: %s", err.Error())
			os.Exit(1)
		}
	case "version":
		err := subcommands.ParseVersionArguments(subCommandArguments)
		if err != nil {
			fmt.Printf("Error: %s", err.Error())
			os.Exit(1)
		}

		err = subcommands.Version()
		if err != nil {
			fmt.Printf("Error: %s", err.Error())
			os.Exit(1)
		}
	default:
		fmt.Printf("Unrecognized sub-command: %s\n", subCommand)
		os.Exit(1)
	}
}

func parseArguments(arguments []string) (string, []string) {
	var command string
	var commandArguments []string

	if len(arguments) == 0 {
		command = "help"
		commandArguments = []string{}
	} else {
		switch arguments[0] {
		case "--help":
			command = "help"
		case "--version":
			command = "version"
		default:
			command = arguments[0]
		}

		commandArguments = arguments[1:]
	}

	return command, commandArguments
}
