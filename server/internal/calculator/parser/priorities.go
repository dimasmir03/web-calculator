package parser

import (
	"errors"
	"math"

	"github.com/dimasmir03/web-calculator-server/internal/calculator/lexer"
)

type TokenAssociativity uint16
type TokenPrecedence uint16

const (
	LeftAssociativity TokenAssociativity = iota
)

type TokenMeta struct {
	Precedence TokenPrecedence

	Associativity TokenAssociativity
}
type TokenPriorities map[lexer.TokenType]TokenMeta

var (
	ErrZeroPrecedenceSet = errors.New("precedence must be greater than 0")
)

func DefaultTokenPriorities() TokenPriorities {

	return TokenPriorities{

		lexer.Addition:     TokenMeta{Precedence: 20},
		lexer.Substraction: TokenMeta{Precedence: 20},

		lexer.Multiplication: TokenMeta{Precedence: 40},
		lexer.Division:       TokenMeta{Precedence: 40},

		lexer.UnaryAddition:     TokenMeta{Precedence: 60},
		lexer.UnarySubstraction: TokenMeta{Precedence: 60},
	}
}

func (tp TokenPriorities) Normalize() error {
	for k := range tp {
		switch k {
		case lexer.Addition, lexer.Substraction, lexer.Multiplication, lexer.Division,
			lexer.UnaryAddition, lexer.UnarySubstraction:

		default:
			delete(tp, k)
		}
	}
	if tp.MinPrecedence() == 0 {
		return ErrZeroPrecedenceSet
	}
	return nil
}

func (tp TokenPriorities) GetMeta(tokenType lexer.TokenType) TokenMeta {
	if meta, ok := tp[tokenType]; ok {
		return meta
	}
	return TokenMeta{}
}

func (tp TokenPriorities) GetPrecedence(tokenType lexer.TokenType) TokenPrecedence {
	return tp.GetMeta(tokenType).Precedence
}

func (tp TokenPriorities) GetAssociativity(tokenType lexer.TokenType) TokenAssociativity {
	return tp.GetMeta(tokenType).Associativity
}

func (tp TokenPriorities) MaxPrecedence() TokenPrecedence {
	maxPrecedence := TokenPrecedence(0)
	for _, v := range tp {
		if v.Precedence > maxPrecedence {
			maxPrecedence = v.Precedence
		}
	}
	return maxPrecedence
}

func (tp TokenPriorities) MinPrecedence() TokenPrecedence {
	minPrecedenc := TokenPrecedence(math.MaxUint16)
	for _, v := range tp {
		if v.Precedence < minPrecedenc {
			minPrecedenc = v.Precedence
		}
	}
	return minPrecedenc
}

func (tp TokenPriorities) NextPrecedence(current TokenPrecedence) TokenPrecedence {
	next := tp.MaxPrecedence()
	for _, v := range tp {
		if v.Precedence > current && v.Precedence < next {
			next = v.Precedence
		}
	}
	return next
}
