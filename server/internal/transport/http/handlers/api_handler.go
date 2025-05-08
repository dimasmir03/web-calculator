package handlers

import (
	"fmt"
	"net/http"

	"github.com/dimasmir03/web-calculator-server/internal/calculator/cmd/calculator"
	"github.com/dimasmir03/web-calculator-server/internal/model"
	"github.com/dimasmir03/web-calculator-server/internal/storage/sqlite"
	"github.com/dimasmir03/web-calculator-server/internal/transport/http/errors"
	"github.com/dimasmir03/web-calculator-server/internal/transport/http/models"
	"github.com/labstack/echo/v4"
)

// GetExpressions godoc
// @Summary Get all expressions
// @Description Get list of all expressions with their statuses
// @Tags expressions
// @Produce json
// @Success 200 {object} models.ExpressionsResponse
// @Router /expressions [get]
func WrapperHandlerGetExpressions(calc *calculator.Calculator) echo.HandlerFunc {
	return func(c echo.Context) error {
		exps := calc.GetExpressionsStatus()
		var res struct {
			Expressions []calculator.Expr `json:"expressions"`
		}
		res.Expressions = exps
		return c.JSON(200, res)
	}
}

// GetExpression godoc
// @Summary Get expression by ID
// @Description Get expression details by ID
// @Tags expressions
// @Produce json
// @Param id path string true "Expression ID"
// @Success 200 {object} models.Expression
// @Failure 404 {object} models.ErrResponse
// @Router /expressions/{id} [get]
func WrapperHandlerGetExpression(calc *calculator.Calculator) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		fmt.Println(id)
		expr := calc.GetExpressionById(id)
		fmt.Println(expr)
		if expr == nil {
			return c.String(http.StatusNotFound, errors.ErrNotFoundID.Error())
		}
		return c.JSON(http.StatusOK, expr)
	}
}

// PostExpression godoc
// @Summary Create new expression
// @Description Add new arithmetic expression for calculation
// @Tags expressions
// @Accept json
// @Produce json
// @Param input body models.CalculateRequest true "Expression data"
// @Success 201 {object} models.CalculateResponse
// @Failure 422 {object} models.ErrResponse
// @Router /calculate [post]
func WrapperHandlerPostExpression(db *sqlite.Storage, calc *calculator.Calculator) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req models.CalculateRequest
		c.Bind(&req)
		id, err := calc.AddExpr(req.Expression)
		if err != nil {
			return c.String(http.StatusUnprocessableEntity, err.Error())
		}
		db.CreateExpression(&model.Expression{ID: id, Expression: req.Expression})
		var res struct {
			Id string `json:"id"`
		}
		res.Id = id
		return c.JSON(http.StatusCreated, res)
	}
}
