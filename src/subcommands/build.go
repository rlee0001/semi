package subcommands

import (
	"errors"
	"flag"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"semi/src/codegen"
	"semi/src/parser"
)

func ParseBuildArguments(arguments []string) (string, string, error) {
	buildCmd := flag.NewFlagSet("build", flag.ExitOnError)
	buildOutput := buildCmd.String("o", "", "file to write the binary to")

	err := buildCmd.Parse(arguments)
	if err != nil {
		return "", "", err
	}

	buildArguments := buildCmd.Args()

	if len(buildArguments) != 1 {
		return "", "", errors.New("command error: <source> argument is required")
	}

	if *buildOutput == "" {
		semiBaseName := filepath.Base(buildArguments[0])
		return buildArguments[0], strings.TrimSuffix(semiBaseName, path.Ext(semiBaseName)), nil
	} else {
		return buildArguments[0], *buildOutput, nil
	}
}

func Build(semiFilePath string, binaryFileName string) error {
	// TODO: Ensure clang/llvm is present.

	module, err := parser.ParseFile(semiFilePath)
	if err != nil {
		return err
	}

	llFilePath, err := codegen.ToLlvmFile(module)
	if err != nil {
		return err
	}

	// TODO: Rename this RunClangBuild? Move that function from codegen to here?
	err = codegen.RunClang(llFilePath, binaryFileName)
	if err != nil {
		return err
	}

	defer func() {
		if err := os.Remove(llFilePath); err != nil {
			log.Fatal(err)
		}
	}()

	return nil
}
