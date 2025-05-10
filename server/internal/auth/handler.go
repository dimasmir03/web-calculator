package auth

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	Service *Service
}

func (h *Handler) Register(с echo.Context) error {
	var req RegisterRequest
	if err := с.Bind(&req); err != nil {
		h.Service.log.WithError(err).Errorf("Error binding register request: %v", err)
		return с.JSON(http.StatusBadRequest, err)
	}
	err := h.Service.Register(req.Login, req.Password)
	if err != nil {
		h.Service.log.WithError(err).Errorf("Error registering user: %v", err)
		return с.JSON(http.StatusInternalServerError, err)
	}
	return с.NoContent(http.StatusOK)
}

func (h *Handler) Login(c echo.Context) error {
	var req LoginRequest

	if err := c.Bind(&req); err != nil {
		h.Service.log.WithError(err).Errorf("Error binding login request: %v", err)
		return c.JSON(http.StatusBadRequest, err)
	}
	token, err := h.Service.Login(req.Login, req.Password)
	if err != nil {
		if err == echo.ErrUnauthorized {
			h.Service.log.WithError(err).Errorf("Invalid login or password: %v", err)
			return c.JSON(http.StatusUnauthorized, "invalid login or password: "+err.Error())
		}
		h.Service.log.WithError(err).Errorf("Error logging in: %v", err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	c.Response().Header().Set("Authorization", token)
	return c.JSON(http.StatusOK, token)
}
