package auth

import (
	"fmt"
	"net/http"

	"github.com/dimasmir03/web-calculator-server/internal/calculator/ast"
	"github.com/labstack/echo/v4"
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

func (h *Handler) bind(c echo.Context, i interface{}) error {
	if err := c.Bind(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("invalid json: %v", err))
	}
	if err := c.Validate(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("invalid fields: %v", err))
	}

	return nil
}
