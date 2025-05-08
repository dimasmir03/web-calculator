package router

import (
	"github.com/dimasmir03/web-calculator-server/internal/calculator/cmd/calculator"
	"github.com/dimasmir03/web-calculator-server/internal/storage/sqlite"
	"github.com/dimasmir03/web-calculator-server/internal/transport/http/handlers"
	"github.com/labstack/echo/v4"
)

func ApiRouter(e *echo.Group, db *sqlite.Storage, calculator *calculator.Calculator) {
	e.Add(echo.GET, "/api/v1/expressions", handlers.WrapperHandlerGetExpressions(calculator))
	e.Add(echo.GET, "/api/v1/expressions/:id", handlers.WrapperHandlerGetExpression(calculator))
	e.Add(echo.POST, "/api/v1/calculate", handlers.WrapperHandlerPostExpression(db, calculator))

}
