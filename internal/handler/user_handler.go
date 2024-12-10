package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

func (h *UserHandler) SignInHandler(c echo.Context) error {
	return c.JSON(http.StatusNoContent, nil)
}
