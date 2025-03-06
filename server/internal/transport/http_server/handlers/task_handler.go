package handlers

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/dimasmir03/web-calculator-server/internal/calculator/cmd/calculator"
	"github.com/dimasmir03/web-calculator-server/internal/transport/http_server/models"
	"github.com/labstack/echo/v4"
)

func WrapperHandlerGetTask(calc *calculator.Calculator) echo.HandlerFunc {
	return func(c echo.Context) error {
		expr := calc.GetSimpleExpr()
		fmt.Println(expr)
		if string(expr.Id) == "" {
			return c.JSON(http.StatusNotFound, "нету задач")
		}
		var task models.TaskResponse
		task.Id = string(expr.Id)
		task.Arg1 = expr.A.(float64)
		task.Arg2 = expr.B.(float64)
		task.Operation = expr.Op
		switch expr.Op {
		case "Addition":
			task.OperationTime, _ = strconv.Atoi(os.Getenv("TIME_ADDITION_MS"))
		case "Substraction":
			task.OperationTime, _ = strconv.Atoi(os.Getenv("TIME_SUBTRACTION_MS"))
		case "Multiplication":
			task.OperationTime, _ = strconv.Atoi(os.Getenv("TIME_MULTIPLICATION_MS"))
		case "Division":
			task.OperationTime, _ = strconv.Atoi(os.Getenv("TIME_DIVISION_MS"))
		}
		return c.JSON(http.StatusOK, task)
	}
}

func WrapperHandlerPostTask(calc *calculator.Calculator) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req models.TaskResultRequest
		err := c.Bind(&req)
		if err != nil {
			return c.JSON(http.StatusUnprocessableEntity, "ошибка декода запроса")
		}
		if err := calc.SetSimpleExprResult(req.Id, req.Result); err != nil {
			return c.JSON(http.StatusNotFound, err.Error())
		}
		return c.NoContent(http.StatusOK)
	}
}
