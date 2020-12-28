package parser

import (
	"errors"
	"path"

	"github.com/antlr/antlr4/runtime/Go/antlr"

	"semi/src/ast"
	"semi/src/parser/gen"
)

// TODO: Add ParseString()?

func ParseFile(semiFilePath string) (ast.Node, error) {
	ext := path.Ext(semiFilePath)

	if ext != ".semi" {
		return nil, errors.New("expected a .semi file for input")
	}

	sourceFileStream, err := antlr.NewFileStream(semiFilePath)
	if err != nil {
		return nil, err
	}

	semiLexer := gen.NewSemiLexer(sourceFileStream)
	semiParser := gen.NewSemiParser(antlr.NewCommonTokenStream(semiLexer, antlr.TokenDefaultChannel))
	semiListener := NewModuleListener()

	antlr.ParseTreeWalkerDefault.Walk(semiListener, semiParser.Module())

	return semiListener.PopNode(), nil
}
