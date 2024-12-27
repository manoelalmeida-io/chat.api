package handler

import (
	"chat_api/internal/model"
	"chat_api/internal/repository"
	"database/sql"
	"net/http"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	userRepository        *repository.UserRepository
	userContactRepository *repository.UserContactRepository
}

func NewUserHandler(
	userRepository *repository.UserRepository,
	userContactRepository *repository.UserContactRepository,
) *UserHandler {
	return &UserHandler{
		userRepository:        userRepository,
		userContactRepository: userContactRepository,
	}
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

func (h *UserHandler) AddContactHandler(c echo.Context) error {
	user := c.Get("userInfo").(*model.User)
	userContact := new(model.UserContact)

	userContact.UserId = &user.Id

	if err := c.Bind(userContact); err != nil {
		return err
	}

	userContact, err := h.userContactRepository.Save(*userContact)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, userContact)
}

func (h *UserHandler) GetContactByIdHandler(c echo.Context) error {
	strIndex := c.Param("id")

	id, err := strconv.Atoi(strIndex)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "parameter id cannot be converted to integer")
	}

	userContact, err := h.userContactRepository.FindById(int64(id))
	if err == sql.ErrNoRows {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, userContact)
}

func (h *UserHandler) UpdateContactHandler(c echo.Context) error {
	strIndex := c.Param("id")

	id, err := strconv.Atoi(strIndex)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	userContact, err := h.userContactRepository.FindById(int64(id))
	if err == sql.ErrNoRows {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	requestUserContact := new(model.UserContact)
	if err := c.Bind(requestUserContact); err != nil {
		return err
	}

	userContact.FirstName = requestUserContact.FirstName
	userContact.LastName = requestUserContact.LastName
	userContact.Email = requestUserContact.Email

	updatedUserContact, err := h.userContactRepository.Save(*userContact)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, updatedUserContact)
}

func (h *UserHandler) DeleteContactHandler(c echo.Context) error {
	strIndex := c.Param("id")

	id, err := strconv.Atoi(strIndex)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "parameter id cannot be converted to integer")
	}

	_, err = h.userContactRepository.FindById(int64(id))
	if err == sql.ErrNoRows {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	if err = h.userContactRepository.DeleteById(int64(id)); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusNoContent, nil)
}

func (h *UserHandler) FindContactsByUserHandler(c echo.Context) error {
	userInfo := c.Get("userInfo").(*model.User)

	userContacts, err := h.userContactRepository.FindByUserId(userInfo.Id)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, userContacts)
}
