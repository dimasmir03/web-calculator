package router

import (
	"github.com/dimasmir03/web-calculator-server/internal/calculator/cmd/calculator"
	"github.com/dimasmir03/web-calculator-server/internal/transport/http/handlers"
	"github.com/labstack/echo/v4"
)

func InternalRouter(e *echo.Echo, calc *calculator.Calculator) {

	e.Add(echo.GET, "/internal/task", handlers.WrapperHandlerGetTask(calc))
	e.Add(echo.POST, "/internal/task", handlers.WrapperHandlerPostTask(calc))
}
