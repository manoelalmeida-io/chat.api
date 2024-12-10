package main

import (
	"chat_api/internal/handler"
	"chat_api/internal/jwt"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Use(echojwt.WithConfig(echojwt.Config{
		KeyFunc:       jwt.KeyFunc,
		SigningMethod: "RS256",
		TokenLookup:   "header:Authorization:Bearer ",
	}))

	userHandler := handler.NewUserHandler()

	e.POST("/users/sign-in", userHandler.SignInHandler)

	e.Logger.Fatal(e.Start(":8080"))
}
