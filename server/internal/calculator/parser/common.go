package parser

import (
	"github.com/dimasmir03/web-calculator-server/internal/calculator/ast"
	"github.com/dimasmir03/web-calculator-server/internal/calculator/lexer"
)

type Parser interface {
	Parse(tokenList []*lexer.Token) (ast.Node, error)
}
