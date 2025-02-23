package router

import (
	"github.com/dimasmir03/web-calculator-server/internal/transport/http_server/handlers"
	"github.com/dimasmir03/web-calculator-server/pkg/calculator/cmd/calculator"
	"github.com/labstack/echo/v4"
)

func ApiRouter(e *echo.Echo, calculator *calculator.Calculator) {
	e.Add(echo.GET, "/api/v1/expressions", handlers.WrapperHandlerGetExpressions(calculator))
	e.Add(echo.GET, "/api/v1/expressions/:id", handlers.WrapperHandlerGetExpression(calculator))
	e.Add(echo.POST, "/api/v1/calculate", handlers.WrapperHandlerPostExpression(calculator))

}
