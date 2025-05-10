package auth

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func (s *Service) JWTMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get(echo.HeaderAuthorization)
			if authHeader == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "missing authorization header")
			}

			// Формат: "Bearer <token>"
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				return echo.NewHTTPError(http.StatusUnauthorized, fmt.Sprintf("invalid authorization format: %s", authHeader))
			}

			token, err := jwt.Parse(parts[1], func(token *jwt.Token) (interface{}, error) {
				return []byte(s.jwtSecret), nil
			})

			if err != nil {
				c.Logger().Errorf("error parsing token: %v", err)
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
			}

			if !token.Valid {
				c.Logger().Errorf("invalid token: %v", token)
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				c.Logger().Errorf("token claims are not a jwt.MapClaims: %T", token.Claims)
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
			}

			userID, ok := claims["sub"]
			if !ok {
				c.Logger().Errorf("token claims do not contain a 'sub' key")
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
			}

			c.Set("userID", userID)
			return next(c)
		}
	}
}
