package subcommands

import (
	"errors"
	"flag"
	"semi/src/parser"
	"semi/src/parsetree"
)

func ParseParseTreeArguments(arguments []string) (string, string, error) {
	parseTreeCmd := flag.NewFlagSet("parse-tree", flag.ExitOnError)
	parseTreeOutput := parseTreeCmd.String("o", "", "file to write the parse tree to")

	err := parseTreeCmd.Parse(arguments)
	if err != nil {
		return "", "", err
	}

	parseTreeArguments := parseTreeCmd.Args()

	if len(parseTreeArguments) != 1 {
		return "", "", errors.New("command error: <source> argument is required")
	}

	return parseTreeArguments[0], *parseTreeOutput, nil
}

func ParseTree(semiFilePath string, outputFilePath string) error {
	// TODO: Print or write parse tree. (Separate functions?)
	// TODO: Fix XML output formatting (xml directive)
	module, err := parser.ParseFile(semiFilePath)
	if err != nil {
		return err
	}

	parseTreeWalker := parsetree.NewParseTreeWalker()

	_, err = parseTreeWalker.VisitNode(module, nil)

	return err
}
