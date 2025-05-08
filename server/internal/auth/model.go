package auth

import (
	"github.com/dimasmir03/web-calculator-server/internal/calculator/ast"
)

type Expr struct {
	Id         ast.UID `json:"id"`
	Expression string  `json:"expression"`
	Status     string  `json:"status"`
	Result     float64 `json:"result"`
}

type RegisterRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
