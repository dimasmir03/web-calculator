package router

import (
	"github.com/dimasmir03/web-calculator-server/internal/transport/http_server/handlers"
	"github.com/dimasmir03/web-calculator-server/pkg/calculator/cmd/calculator"
	"github.com/labstack/echo/v4"
)

func InternalRouter(e *echo.Echo, calc *calculator.Calculator) {

	e.Add(echo.GET, "/internal/task", handlers.WrapperHandlerGetTask(calc))
	e.Add(echo.POST, "/internal/task", handlers.WrapperHandlerPostTask(calc))
}
