package module

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type User struct{}

func (ref User) Route(e *echo.Group) {
	handler := UserHandler{}

	e.GET("/:project_key/example-trigger/:value", handler.Trigger)

}

// ---------------------------------------------------------------------------------------------
// ---------------------------------------------------------------------------------------------

type UserHandler struct{}

func (handler UserHandler) Trigger(c echo.Context) error {
	// var err error

	value := c.Param("value")

	return c.JSON(http.StatusOK, map[string]string{"message": "OK", "value": value})
}
