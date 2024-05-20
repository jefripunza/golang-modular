package module

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Cart struct{}

func (ref Cart) Route(e *echo.Group) {
	handler := CartHandler{}

	e.GET("/:project_key/example-trigger/:value", handler.Trigger)

}

// ---------------------------------------------------------------------------------------------
// ---------------------------------------------------------------------------------------------

type CartHandler struct{}

func (handler CartHandler) Trigger(c echo.Context) error {
	// var err error

	value := c.Param("value")

	return c.JSON(http.StatusOK, map[string]string{"message": "OK", "value": value})
}
