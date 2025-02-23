package middleware

import (
	"time"

	"github.com/labstack/echo/v4"
)

// middleware для логирования
func logging(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		start := time.Now()
		err := next(c)
		if err != nil {
			return err
		}
		c.Logger().Infof("HTTP Запрос %s %s %s", c.Request().Method, c.Request().RequestURI, time.Since(start))
		return nil
	}
}
