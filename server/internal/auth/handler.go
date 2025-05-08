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
		return с.JSON(http.StatusBadRequest, err)
	}
	err := h.Service.Register(req.Login, req.Password)
	if err != nil {
		return с.JSON(http.StatusInternalServerError, err)
	}
	return с.NoContent(http.StatusOK)
}

func (h *Handler) Login(c echo.Context) error {
	var req LoginRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	token, err := h.Service.Login(req.Login, req.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	c.Response().Header().Set("Authorization", token)
	return c.JSON(http.StatusOK, token)
}
