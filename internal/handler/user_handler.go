package handler

import (
	"chat_api/internal/model"
	"chat_api/internal/repository"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	userRepository *repository.UserRepository
}

func NewUserHandler(userRepository *repository.UserRepository) *UserHandler {
	return &UserHandler{userRepository: userRepository}
}

func (h *UserHandler) SignInHandler(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	firstName := claims["given_name"].(string)
	lastName := claims["family_name"].(string)
	email := claims["email"].(string)
	googleSub := claims["sub"].(string)

	newUser := &model.User{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		GoogleSub: googleSub,
	}

	existent, err := h.userRepository.FindBySub(googleSub)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if existent != nil {
		return echo.NewHTTPError(http.StatusConflict, "user already exists")
	}

	newUser, err = h.userRepository.Save(*newUser)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, newUser)
}
