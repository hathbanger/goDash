package server

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

func accessible(c echo.Context) error {

	return c.String(http.StatusOK, "Hello from accessible!")
}

func restricted(c echo.Context) error {

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	id := claims["id"].(string)

	return c.String(http.StatusOK, "Welcome "+id+"!")
}
