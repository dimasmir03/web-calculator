package evaluator

import (
	"errors"
	"fmt"

	ast2 "github.com/dimasmir03/web-calculator-server/internal/calculator/ast"
)

type FunctionHandler struct {
	Description  string
	Handler      func(x ...float64) (float64, error)
	MinArguments int
	MaxArguments int
	ArgsNames    []string
}

type NumericEvaluator struct {
}

type VariableTuple struct {
	Name  string
	Value float64
}
type FunctionTuple struct {
	Name     string
	Function FunctionHandler
}

func NewNumericEvaluator() (*NumericEvaluator, error) {

	return &NumericEvaluator{}, nil
}

func (e *NumericEvaluator) Eval(rootNode ast2.Node) (float64, error) {
	switch n := rootNode.(type) {
	case *ast2.BinaryNode:
		return e.handleBinary(n)
	case *ast2.UnaryNode:
		return e.handleUnary(n)
	case *ast2.NumericNode:
		return n.Value(), nil
	}
	return 0, EvalError(rootNode.GetToken(), fmt.Errorf("unimplemented node type %T", rootNode))
}

func (e *NumericEvaluator) handleBinary(n *ast2.BinaryNode) (float64, error) {
	l, err := e.Eval(n.Left())
	if err != nil {
		return 0, err
	}
	r, err := e.Eval(n.Right())
	if err != nil {
		return 0, err
	}

	switch n.Operator() {
	case ast2.Addition:
		return l + r, nil
	case ast2.Substraction:
		return l - r, nil
	case ast2.Multiplication:
		return l * r, nil
	case ast2.Division:
		return l / r, nil
	}

	return 0, EvalError(n.GetToken(), fmt.Errorf("unimplemented operator %s", n.Operator()))
}

func (e *NumericEvaluator) handleUnary(n *ast2.UnaryNode) (float64, error) {
	val, err := e.Eval(n.Next())
	if err != nil {
		return 0, err
	}

	switch n.Operator() {
	case ast2.Substraction:
		return -val, nil
	case ast2.Addition:
		return val, nil
	}

	return 0, EvalError(n.GetToken(), errors.New("unary node supports only Addition and Substraction operator"))
}
