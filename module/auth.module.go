package module

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Auth struct{}

func (ref Auth) Route(api *echo.Group) {
	handler := AuthHandler{}
	route := api.Group("/auth")

	route.GET("/login", handler.Login)

}

// ---------------------------------------------------------------------------------------------
// ---------------------------------------------------------------------------------------------

type AuthHandler struct{}

func (handler AuthHandler) Login(c echo.Context) error {
	// var err error

	return c.JSON(http.StatusOK, map[string]string{"message": "OK"})
}
